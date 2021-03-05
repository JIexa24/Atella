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
	"../../../atella/logging"
	"../../../atella/reporter"
)

// Config structure description configuration for start service.
type Config struct {
	reporterWorker *reporter.Reporter
	finalize       chan struct{}
	printVersion   bool
	configFilePath string
	configDirPath  string
	GitCommit      string
	GoVersion      string
	Version        string
	Service        string
	BinPrefix      string
	ScriptsPrefix  string
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
	finalize:       make(chan struct{}),
	configFilePath: "",
	configDirPath:  "",
	printVersion:   false,
	GitCommit:      GitCommit,
	GoVersion:      GoVersion,
	Version:        Version,
	Service:        Service}
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
	configuration.ReadConfig(global.configFilePath, global.configDirPath, config)

	// set hostname by os.Hostname if it not defined.
	if config.Hostname == "" {
		config.Hostname, err = os.Hostname()
		if err != nil {
			logger.Fatalf("can't get hostname [%s]", err.Error())
			os.Exit(1)
		}
	}

	logger, err = logging.ConfigureLog(
		config.Logger.LogFile, config.Logger.LogLevel, global.Service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not configure log: %s\n", err.Error())
		os.Exit(1)
	}

	defer logger.Infof("%s stopped. Version: %s", global.Service, global.Version)
	r, err := reporter.Worker(config.Reporter, config.Hostname,
		config.Channels, logger)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	global.reporterWorker = r
	global.reporterWorker.Start()

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
