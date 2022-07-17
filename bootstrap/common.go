package bootstrap

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
	"zerologix-homework/src/pkg/util"

	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_WORK_DIR  = "ZerologixHomework"
	DEFAULT_IANA_ZONE = "Asia/Taipei"
	DEFAULT_CONFIG    = "test"
)

var (
	lock sync.RWMutex
	cfg  *Config
)

func Get() (*Config, error) {
	lock.RLock()
	isNoCfg := cfg == nil
	lock.RUnlock()
	if isNoCfg {
		lock.Lock()
		defer lock.Unlock()
		if cfg == nil {
			err := loadConfig()
			if err != nil {
				return nil, err
			}
		}
	}
	copy := *cfg
	return &copy, nil
}

func loadConfig() error {
	configName := GetEnvConfig()
	if configName == "" {
		configName = DEFAULT_CONFIG
	}

	var err error
	cfg, err = loadYmlConfig(configName)
	if err != nil {
		return err
	}

	loadDefault()

	if err := loadEnv(); err != nil {
		return err
	}

	return nil
}

func loadYmlConfig(fileName string) (*Config, error) {
	root, err := GetRootDirPath()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join("config")
	path := fmt.Sprintf("%s/%s/%s.yml", root, configDir, fileName)

	cfgBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(cfgBytes, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadDefault() {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.Var.TimeZone == "" {
		cfg.Var.TimeZone = DEFAULT_IANA_ZONE
	}
}

func loadEnv() error {
	if cfg == nil {
		cfg = &Config{}
	}

	if envStr := GetEnvPort(); envStr != "" {
		port, err := strconv.Atoi(envStr)
		if err != nil {
			return err
		}
		cfg.Server.Port = port
	}

	return nil
}

func GetRootDirPath() (string, error) {
	if dir := GetEnvWorkDir(); dir != "" {
		return util.GetRootOf(dir)
	}

	return ".", nil
}
