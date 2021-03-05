package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"../../../atella"
	"../../../atella/configuration"
	"../../../atella/httpclient"
	"../../../atella/httpserver"
	"../../../atella/logging"
	"../../../atella/reporter"
)

// Config structure description configuration for start service.
type Config struct {
	client *httpclient.Client
	server *httpserver.Server
	// waitGroup      *sync.WaitGroup
	reporterWorker  *reporter.Reporter
	finalize        chan struct{}
	printVersion    bool
	withoutReporter bool
	configFilePath  string
	configDirPath   string
	GitCommit       string
	GoVersion       string
	Version         string
	Service         string
	BinPrefix       string
	ScriptsPrefix   string
}

// Vars, linked by Makefile.
var (
	GitCommit string = "unknown"
	GoVersion string = "unknown"
	Version   string = "unknown"
	Service   string = "Atella"
)

// Initializing configuration.
var global Config = Config{
	client:          nil,
	server:          nil,
	finalize:        make(chan struct{}),
	configFilePath:  "",
	configDirPath:   "",
	printVersion:    false,
	withoutReporter: false,
	GitCommit:       GitCommit,
	GoVersion:       GoVersion,
	Version:         Version,
	Service:         Service}
var logger atella.Logger

// handle is interrupts handler.
func handle(c chan os.Signal) {
	for {
		sig := <-c
		logger.Infof("Receive %s [%s]",
			sig, sig.String())
		switch sig.String() {
		case "hangup":

		case "interrupt":
			if !global.withoutReporter {
				global.reporterWorker.StopReporter()
			}
			close(global.finalize)
		case "user defined signal 1":
		case "user defined signal 2":
		default:
		}
	}
}

// usage is function is a handler for runtime flag -h.
func usage() {
	fmt.Fprintf(os.Stderr, "[%s] Usage: %s [params]\n", global.Service,
		os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

// parseFlags is function initialize application runtime flags.
func parseFlags() {
  arch := runtime.GOARCH
	flag.Usage = usage
	flag.StringVar(&global.configFilePath, "config", "/etc/atella/atella.yml",
		"Path to config")
	flag.StringVar(&global.configDirPath, "config-dir", "/etc/atella/conf.d",
		"Path to config dir")
	flag.BoolVar(&global.printVersion, "version", false, "Print version and exit")
	flag.BoolVar(&global.withoutReporter, "without-reporter", false,
		"Don.t start reporter")
	flag.Parse()

	if global.printVersion {
		fmt.Println("Atella")
		fmt.Println("Version:", global.Version)
		fmt.Println("Arch:", arch)
		fmt.Println("Git Commit:", global.GitCommit)
		fmt.Println("Go Version:", global.GoVersion)
		os.Exit(0)
	}
}

// main is main ^).
func main() {
	var err error = nil
	parseFlags()

	config := configuration.GetDefault()
	configuration.PrintConfig(config)
	// if *printDefaultConfigFlag {
	// 	cmd.PrintConfig(config)
	// 	os.Exit(0)
	// }
	configuration.ReadConfig(global.configFilePath, global.configDirPath, config)
	// set hostname by os.Hostname if it not defined.
	if config.Hostname == "" {
		config.Hostname, err = os.Hostname()
		if err != nil {
			logger.Fatalf("can't get hostname [%s]", err.Error())
			os.Exit(1)
		}
	}
	configuration.PrintConfig(config)
	logger, err = logging.ConfigureLog(
		config.Logger.LogFile, config.Logger.LogLevel, global.Service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not configure log: %s\n", err.Error())
		os.Exit(1)
	}
	parseResult, err := config.ParseHosts(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not parse hosts: %s\n", err.Error())
		os.Exit(1)
	}

	logger.Infof("Self indexes %v", parseResult.SelfIndexes)
	logger.Infof("Master servers %v", parseResult.MasterIndexes)
	logger.Infof("Sectors %#v", parseResult.SectorMapper)

	defer logger.Infof("%s stopped. Version: %s", global.Service, global.Version)
	if !global.withoutReporter {
		r, err := reporter.Worker(config.Reporter, config.Hostname,
			config.Channels, logger)
		if err != nil {
			logger.Fatal(err)
			os.Exit(1)
		}

		global.reporterWorker = r
		global.reporterWorker.Start()
	}

	var vectorChannel chan map[string]atella.HostVector = nil
	global.client, vectorChannel = httpclient.NewHTTPClient(
		config.Hostname, config.Security,
		config.Hosts, parseResult,
		config.Connectivity, logger)
	global.client.Start()
	global.server = httpserver.NewHTTPServer(config.Hosts, parseResult,
		global.Version, config.Hostname, config.Security, logger)
	global.server.Start(vectorChannel)

	// // Creating signals handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGUSR1)
	signal.Notify(c, syscall.SIGUSR2)

	go handle(c)

	logger.Infof("Started %s version %s", global.Service, global.Version)
	<-global.finalize
}
