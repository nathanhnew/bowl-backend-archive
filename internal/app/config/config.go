package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

var DefaultConfigLocation string = "config/main.conf.json"

type Config struct {
	values   map[string]interface{}
	location string
}

func GetConfig(location string) (*Config, error) {
	if location == "" {
		location = DefaultConfigLocation
	}
	cfgFile, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()
	byteCfg, _ := ioutil.ReadAll(cfgFile)

	var values map[string]interface{}
	json.Unmarshal([]byte(byteCfg), &values)

	var cfg Config
	cfg.values = values
	cfg.location = location

	return &cfg, nil
}

func (config *Config) GetListenPort() int {
	return int(config.values["port"].(float64))
}

func (config *Config) GetMongoUri() string {
	return config.values["mongoUri"].(string)
}

func (config *Config) GetMongoTimeout() time.Duration {
	return config.values["timeout"].(time.Duration)
}
