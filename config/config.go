package config

import (
	"encoding/json"
	"os"
)

const DefaultConfigFilename = "config.json"

type Config struct {
	Database struct {
		Path string `json:"path"`
	} `json:"database"`
	Http struct {
		Addr string `json:"addr"`
	} `json:"http"`
	TGBot struct {
		Token string `json:"token"`
	} `json:"tgbot"`
}

func ReadConfig(path string) (*Config, error) {
	name := DefaultConfigFilename
	if path != "" {
		name = path
	}

	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
