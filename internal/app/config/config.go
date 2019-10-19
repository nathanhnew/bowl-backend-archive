package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

var DefaultConfigLocation string = "config/main.conf.json"

type Config struct {
	Values   map[string]interface{}
	Location string
}

func GetConfig(location string) (*Config, error) {
	var values map[string]interface{}
	var cfg Config
	if location == "" {
		location = DefaultConfigLocation
	}
	cfgFile, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()
	byteCfg, _ := ioutil.ReadAll(cfgFile)

	json.Unmarshal([]byte(byteCfg), &values)

	cfg.Values = values
	cfg.Location = location

	return &cfg, nil
}

func (config *Config) GetListenPort() int {
	if port, ok := config.Values["port"]; ok {
		return int(port.(float64))
	}
	return 0
}

func (config *Config) GetMongoUri() string {
	if uri, ok := config.Values["mongoUri"]; ok {
		return uri.(string)
	}
	return ""
}

func (config *Config) GetMongoTimeout() time.Duration {
	if timeout, ok := config.Values["timeout"]; ok {
		return timeout.(time.Duration)
	}
	return time.Second
}
