package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"../../atella"
)

// GetDefault return default config structure.
func GetDefault() *atella.Config {
	local := &atella.Config{
		Hostname: "",
		Reporter: atella.ReporterConfig{
			HexLen:      10,
			MessagePath: "/usr/share/atella/msg",
		},
		Logger: atella.LoggerConfig{
			LogFile:  "stdout",
			LogLevel: "info",
		},
	}
	return local
}

// PrintConfig prints config to stdout.
func PrintConfig(config interface{}) {
	d, _ := yaml.Marshal(&config)
	fmt.Println(string(d))
}

// ReadConfig parses config file by the given path.
func ReadConfig(configFileName, configDirName string, config interface{}) error {
	configYaml, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return fmt.Errorf("can't read file [%s] [%s]", configFileName, err.Error())
	}
	// Read config directory
	dirContent, err := readDir(configDirName)
	if err != nil {
		return fmt.Errorf("can't parse config dir [%s] [%s]", configDirName, err.Error())
	}
	// Separate configs by \n because it's required by yaml format.
	configYaml = append(configYaml, '\n')
	configYaml = append(configYaml, dirContent...)
	err = yaml.Unmarshal(configYaml, config)
	if err != nil {
		return fmt.Errorf("can't parse config file [%s] [%s]", configFileName, err.Error())
	}
	return nil
}

// readDir parses config files by the given directory path.
func readDir(configDirName string) ([]byte, error) {
	var dirContent []byte = make([]byte, 0)
	if configDirName == "" {
		return dirContent, nil
	}
	// define function for directory walking
	walkfn := func(thispath string, info os.FileInfo, _ error) error {
		if info == nil {
			return fmt.Errorf("I don't have permissions to read %s",
				thispath)
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), "..") {
				return filepath.SkipDir
			}
			return nil
		}

		name := info.Name()
		if len(name) < 5 || name[len(name)-4:] != ".yml" {
			return nil
		}
		content, err := ioutil.ReadFile(thispath)
		if err != nil {
			return err
		}
		dirContent = append(dirContent, '\n')
		dirContent = append(dirContent, content...)
		return nil
	}

	return dirContent, filepath.Walk(configDirName, walkfn)
}
