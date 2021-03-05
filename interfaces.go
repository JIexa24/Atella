package atella

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Logger implements logger abstraction.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
}

// Config is configuration structure for atella.
type Config struct {
	Hostname     string         `yaml:"hostname"`
	Logger       LoggerConfig   `yaml:"log"`
	PidFile      string         `yaml:"pid_file"`
	ProcFile     string         `yaml:"proc_file"`
	Connectivity int64          `yaml:"connectivity"`
	Reporter     ReporterConfig `yaml:"reporter"`
	Master       bool           `yaml:"master"`
	Interval     string         `yaml:"interval"`
	NetTimeout   string         `yaml:"net_timeout"`
	Security     string         `yaml:"security"`
	// Sectors SectorsConfig `yaml:"sectors"`
	// Masters MastersConfig `yaml:"masters"`
	Channels []map[string]string `yaml:"channels"`
	Hosts    []Host              `yaml:"hosts"`
}

// ReporterConfig is configure log.
type ReporterConfig struct {
	MessagePath string `yaml:"message_path"`
	HexLen      int64  `yaml:"hex_len"`
}

// LoggerConfig is configure log.
type LoggerConfig struct {
	LogFile  string `yaml:"log_file"`
	LogLevel string `yaml:"log_level"`
}

// Host describe every host in config.
type Host struct {
	Address  string   `yaml:"address"`
	Port     string   `yaml:"port"`
	Hostname string   `yaml:"hostname"`
	Web      bool     `yaml:"web"`
	Master   bool     `yaml:"master"`
	Sectors  []string `yaml:"sectors"`
}

// Sector is struct, which describes a sector.
// Contains self indexes and host indexes.
type Sector struct {
	SelfIndexes  []int64
	HostsIndexes []int64
}

// ParsedHosts is a struct for parsing hosts. Contains self hosts indexes,
// masters and sectors mapper.
type ParsedHosts struct {
	SelfIndexes   []int64
	MasterIndexes []int64
	SectorMapper  map[string]Sector
}

// ParseHosts parsing hosts and check duplicates address:port. Also return
// self indexes in array and sectors map.
func (config *Config) ParseHosts(logger Logger) (*ParsedHosts, error) {
	var parseResult *ParsedHosts = &ParsedHosts{
		SelfIndexes:   make([]int64, 0),
		MasterIndexes: make([]int64, 0),
		SectorMapper:  make(map[string]Sector),
	}
	var hostMapper map[string]string = make(map[string]string)
	var duplicateNamesMapper map[string]int64 = make(map[string]int64)

	for index, host := range config.Hosts {
		if host.Master {
			parseResult.MasterIndexes = append(parseResult.MasterIndexes, int64(index))
		}
		if _, ok := duplicateNamesMapper[host.Hostname]; ok {
			duplicateNamesMapper[host.Hostname]++
		} else {
			duplicateNamesMapper[host.Hostname] = 1
		}
		if _, ok := hostMapper[host.Address]; ok {
			if host.Port == hostMapper[host.Address] {
				return nil, fmt.Errorf("duplicate %v:%v", host.Address, host.Port)
			}
		} else {
			hostMapper[host.Address] = host.Port
		}
		isMe := false
		if host.Hostname == config.Hostname {
			isMe = true
			parseResult.SelfIndexes = append(parseResult.SelfIndexes, int64(index))
		}
		for _, s := range host.Sectors {
			if _, ok := parseResult.SectorMapper[s]; !ok {
				parseResult.SectorMapper[s] = Sector{
					HostsIndexes: make([]int64, 0),
					SelfIndexes:  make([]int64, 0),
				}
			}
			sector := parseResult.SectorMapper[s]
			sector.HostsIndexes = append(sector.HostsIndexes, int64(index))
			if isMe {
				sector.SelfIndexes = append(sector.SelfIndexes, int64(index))
			}
			parseResult.SectorMapper[s] = sector
		}
	}
	if config.Connectivity > 0 {
		for _, count := range duplicateNamesMapper {
			if count > config.Connectivity {
				config.Connectivity = count
			}
		}
		logger.Infof("Set connectivity to %d", config.Connectivity)
	}
	if len(parseResult.MasterIndexes) == 0 {
		parseResult.MasterIndexes = append(parseResult.MasterIndexes, -1)
	}
	return parseResult, nil
}

// HostVector is vector element structure.
type HostVector struct {
	Address   string `json:"address"`
	Port      string `json:"port"`
	Hostname  string `json:"hostname"`
	Status    bool   `json:"status"`
	timestamp int64
}

// MasterVector implement states master vector. Has mutex and get/set methods.
type MasterVector struct {
	mutex sync.RWMutex
	List  map[string]map[string]HostVector `json:"list"`
}

// Create new vector and return it.
func NewMasterVector() *MasterVector {
	return &MasterVector{
		mutex: sync.RWMutex{},
		List:  make(map[string]map[string]HostVector),
	}
}

// GetElement get vector in JSON format.
func (v *MasterVector) GetVectorJSON() []byte {
	v.mutex.RLock()
	bytes, _ := json.Marshal(v)
	v.mutex.RUnlock()
	return bytes
}

// SetElement add/update element for host by address:port.
func (v *MasterVector) SetElement(key string, element map[string]HostVector) {
	v.mutex.Lock()
	if v.List[key] == nil {
		v.List[key] = make(map[string]HostVector)
	}
	v.List[key] = element
	v.mutex.Unlock()
}

// Create copy of vector and return it.
func (v *MasterVector) GetVectorCopy() *MasterVector {
	v.mutex.RLock()
	copyV := &MasterVector{
		mutex: sync.RWMutex{},
		List:  v.List,
	}
	v.mutex.RUnlock()
	return copyV
}

// Vector implement states vector. Has mutex and get/set methods.
type Vector struct {
	mutex sync.RWMutex
	List  map[string]HostVector `json:"list"`
}

// GetElement get vector in JSON format.
func (v *Vector) GetVectorJSON() []byte {
	v.mutex.RLock()
	bytes, _ := json.Marshal(v)
	v.mutex.RUnlock()
	return bytes
}

// SetElement add/update element for host by address:port.
func (v *Vector) SetElement(addr, port string, element HostVector) {
	v.mutex.Lock()
	mapAddr := fmt.Sprintf("%s:%s", addr, port)
	element.timestamp = time.Now().Unix()
	v.List[mapAddr] = element
	v.mutex.Unlock()
}

// GetElement get element for host.
func (v *Vector) GetElement(addr, port string) (HostVector, bool) {
	v.mutex.RLock()
	mapAddr := fmt.Sprintf("%s:%s", addr, port)
	var el HostVector
	isNewEl := false
	if _, ok := v.List[mapAddr]; ok {
		el = v.List[mapAddr]
	} else {
		isNewEl = true
		el = NewHostVector()
	}
	v.mutex.RUnlock()
	return el, isNewEl
}

// Create copy of vector and return it.
func (v *Vector) GetVectorCopy() *Vector {
	v.mutex.RLock()
	copyV := &Vector{
		mutex: sync.RWMutex{},
		List:  v.List,
	}
	v.mutex.RUnlock()
	return copyV
}

// Create new vector and return it.
func NewVector() *Vector {
	return &Vector{
		mutex: sync.RWMutex{},
		List:  make(map[string]HostVector),
	}
}

// Create new vector and return it.
func NewHostVector() HostVector {
	return HostVector{
		Address:   "",
		Port:      "",
		Hostname:  "",
		Status:    false,
		timestamp: -1,
	}
}
