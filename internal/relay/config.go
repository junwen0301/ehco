package relay

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type RelayConfig struct {
	Listen        string `json:"listen"`
	ListenType    string `json:"listen_type"`
	Remote        string `json:"remote"`
	TransportType string `json:"transport_type"`
}

type Config struct {
	PATH    string
	Configs []RelayConfig
}

type JsonConfig struct {
	Configs []RelayConfig `json:"configs"`
}

func NewConfig(path string) *Config {
	return &Config{PATH: path, Configs: []RelayConfig{}}
}

func (c *Config) LoadConfig() error {
	if strings.Contains(c.PATH, "http") {
		return c.readFromHttp()
	}
	return c.readFromFile()
}

func (c *Config) readFromFile() error {
	file, err := ioutil.ReadFile(c.PATH)
	if err != nil {
		return err
	}
	jsonConfig := JsonConfig{}
	err = json.Unmarshal([]byte(file), &jsonConfig)
	if err != nil {
		return err
	}
	c.Configs = jsonConfig.Configs
	Logger.Info("load config from file:", c.PATH)
	return nil
}

func (c *Config) readFromHttp() error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(c.PATH)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	jsonConfig := JsonConfig{}
	json.NewDecoder(r.Body).Decode(&jsonConfig)
	c.Configs = jsonConfig.Configs
	Logger.Info("load config from http:", c.PATH, c.Configs)
	return nil
}
