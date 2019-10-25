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

func GetCurrentSeason() string {
	now := time.Now()
	if now.Month() < 0 {
		return string(now.Year() - 1)
	} else {
		return string(now.Format("2006"))
	}
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
		return time.Duration(timeout.(float64))
	}
	return time.Second
}

func (config *Config) GetTokenTimeout() float64 {
	if timeout, ok := config.Values["timeout"]; ok {
		return timeout.(float64)
	}
	return 1
}

func (config *Config) GetValidationKey() []byte {
	if key, ok := config.Values["validationKey"]; ok {
		return []byte(key.(string))
	}
	return []byte("")
}

func (config *Config) GetMongoDB() string {
	if db, ok := config.Values["mongoDatabase"]; ok {
		return db.(string)
	}
	return ""
}
