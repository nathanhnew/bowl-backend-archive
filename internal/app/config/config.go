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
	return int(config.Values["port"].(float64))
}

func (config *Config) GetMongoUri() string {
	return config.Values["mongoUri"].(string)
}

func (config *Config) GetMongoTimeout() time.Duration {
	return config.Values["timeout"].(time.Duration)
}
