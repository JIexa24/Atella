package reporter

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"../../atella"
	"./mail"
	"./tgsibnet"
	// "../reporter/graphite"
)

// Reporter is reporter config structure
type Reporter struct {
	MessagePath string
	hostname    string
	HexLen      int64
	Channels    map[string]interface {
		Send(text string, hostname string) (bool,
			error)
	}
	stopRequest chan struct{}
	waitGroup   *sync.WaitGroup
	logger      atella.Logger
	isLocked    bool
	firstRun    bool
	messageCnt  int
	mutexCnt    sync.Mutex
	mutex       sync.Mutex
}

type msg struct {
	Target  string `json:"target"`
	Message string `json:"message"`
}

var (
	defaultChannels []string = []string{"tgsibnet", "mail"}
)

// RandomHex return a pseudo-random generated hex-string.
// String length are specifyed by config (And get to function as n).
func RandomHex(n int64) (string, error) {
	if n < 0 {
		return "", fmt.Errorf("Length must be grater than 0")
	}
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Worker implement worker logic for reports.
func Worker(r atella.ReporterConfig, hostname string,
	Channels []map[string]string, logger atella.Logger) (*Reporter, error) {
	worker := Reporter{
		MessagePath: r.MessagePath,
		HexLen:      r.HexLen,
		waitGroup:   &sync.WaitGroup{},
		mutexCnt:    sync.Mutex{},
		mutex:       sync.Mutex{},
		stopRequest: make(chan struct{}),
		isLocked:    false,
		firstRun:    true,
		hostname:    hostname,
		messageCnt:  0,
		logger:      logger,
		Channels: make(map[string]interface {
			Send(text string, hostname string) (bool,
				error)
		}),
	}
	if err := worker.parseChannels(Channels); err != nil {
		return nil, err
	}
	return &worker, nil
}

func (worker *Reporter) parseChannels(Channels []map[string]string) error {
	var err error = nil
	var port int = -1
	var en bool = false
	var auth bool = false

	for _, channel := range Channels {
		t := strings.ToLower(channel["type"])
		switch t {
		case "tgsibnet":
			if port, err = strconv.Atoi(channel["port"]); err != nil {
				return err
			}
			if en, err = strconv.ParseBool(channel["enabled"]); err != nil {
				return err
			}
			to := strings.Split(channel["to"], ",")
			worker.Channels[t] = &tgsibnet.AtellaTgSibnetConfig{
				Address:    channel["address"],
				Port:       int16(port),
				Protocol:   channel["protocol"],
				To:         to,
				Enabled:    en,
				NetTimeout: int(2),
			}
		case "mail":
			if port, err = strconv.Atoi(channel["port"]); err != nil {
				return err
			}
			if en, err = strconv.ParseBool(channel["enabled"]); err != nil {
				return err
			}
			if auth, err = strconv.ParseBool(channel["auth"]); err != nil {
				return err
			}
			re := regexp.MustCompile(`@hostname$`)
			channel["from"] = re.ReplaceAllString(
				channel["from"],
				fmt.Sprintf("@%s", worker.hostname))
			to := strings.Split(channel["to"], ",")
			worker.Channels[t] = &mail.AtellaMailConfig{
				Address:    channel["address"],
				Port:       int16(port),
				Auth:       auth,
				Username:   channel["username"],
				Password:   channel["password"],
				From:       channel["from"],
				To:         to,
				Enabled:    en,
				NetTimeout: int(2),
			}
		case "graphite":
		default:
			return fmt.Errorf("unknown channel type [%s]", channel["type"])
		}
	}
	if err != nil {
		return err
	}
	return nil
}

// Start reporter worker.
func (worker *Reporter) Start() {
	worker.waitGroup.Add(1)
	go worker.reporter()
}

// StopReporter closed stopRequest channel and stop reporter.
func (worker *Reporter) StopReporter() {
	worker.logger.Info("Stopping reporter")
	close(worker.stopRequest)
	worker.waitGroup.Wait()
}

func (worker *Reporter) reporter() {
	defer worker.waitGroup.Done()
	for {
		timer := time.NewTimer(time.Second * 10)
		select {
		case <-worker.stopRequest:
			timer.Stop()
			worker.logger.Info("Reporter stopped")
			return
		case <-timer.C:
			worker.Send()
		}
	}
}

// Send implement send-report mechanism. Use files created by Report function.
func (worker *Reporter) Send() {
	var (
		message string = ""
		target  string = ""
		res     bool   = true
		m       msg
	)
	if worker.isLocked {
		worker.logger.Info("Sender iteration already in progress")
		return
	}
	worker.mutex.Lock()
	defer worker.mutex.Unlock()
	worker.isLocked = true
	worker.logger.Infof("Tick reporter iteration [firstRun - %v| Cnt - %v]",
		worker.firstRun, worker.messageCnt)
	if worker.firstRun || worker.messageCnt > 0 {
		files, err := ioutil.ReadDir(worker.MessagePath)
		if err != nil {
			worker.logger.Error(err.Error())
		}
		worker.firstRun = false
		filesCnt := 0
		for _, file := range files {
			if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
				f, err := os.Open(fmt.Sprintf("%s/%s", worker.MessagePath,
					file.Name()))
				if err != nil {
					worker.logger.Error(err.Error())
					continue
				}
				data := make([]byte, file.Size())

				for {
					n, err := f.Read(data)
					if err == io.EOF {
						break
					}
					message = message + string(data[:n])
				}
				f.Close()
				filesCnt = filesCnt + 1
				if err = json.Unmarshal(data, &m); err != nil {
					worker.logger.Errorf(err.Error())
				}
				worker.logger.Infof("Read msg - %s [msg: %s|target: %s]",
					file.Name(), m.Message, m.Target)

				target = strings.ToLower(m.Target)
				if worker.Channels[target] != nil {
					if res, err = worker.Channels[target].Send(
						m.Message, worker.hostname); err != nil {
						worker.logger.Errorf(err.Error())
					}
				} else {
					worker.logger.Errorf("Unsopported channel - %s", target)
				}

				if res {
					worker.mutexCnt.Lock()
					if worker.messageCnt > 0 {
						worker.messageCnt = worker.messageCnt - 1
					}
					worker.mutexCnt.Unlock()
					os.Remove(fmt.Sprintf("%s/%s", worker.MessagePath, file.Name()))
				}
			}
		}
		// conf.SendMetric_AtellaSender("queue.cnt", float64(filesCnt))
	}
	worker.isLocked = false
}

// Report save report as a file (filename are random hex string).
func (worker *Reporter) Report(message, target string) string {
	var (
		hash string = ""
		path string = ""
		m    *msg   = &msg{
			Message: "",
			Target:  ""}
		file    *os.File = nil
		err     error    = nil
		targets []string = make([]string, 0)
	)
	if strings.ToLower(target) == "all" {
		targets = defaultChannels
	} else {
		targets = append(targets, target)
	}

	for _, target := range targets {
		for {
			hash, _ = RandomHex(worker.HexLen)
			path = fmt.Sprintf("%s/%s", worker.MessagePath, hash)
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				break
			}
		}
		file, err = os.Create(path)
		if err != nil {
			worker.logger.Errorf("Unable to create file: %s", err.Error())
		}

		defer file.Close()
		m.Message = message
		m.Target = target
		js, _ := json.Marshal(m)
		file.Write([]byte(js))
		worker.logger.Infof("File - %s [msg: %s|target: %s]",
			path, message, target)
	}

	worker.mutexCnt.Lock()
	worker.messageCnt = worker.messageCnt + 1
	worker.mutexCnt.Unlock()
	return hash
}
