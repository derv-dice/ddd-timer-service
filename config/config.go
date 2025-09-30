package config

import (
	"encoding/json"
	"os"
)

const DefaultConfigFilename = "config.json"

type Config struct {
	LogLevel string `json:"log_level"`
	Database struct {
		Path string `json:"path"`
	} `json:"database"`
	Http struct {
		Addr string `json:"addr"`
	} `json:"http"`
	TGBot struct {
		Token string `json:"token"`
	} `json:"tgbot"`
	Limits struct {
		CalendarImg struct {
			CacheSizeMB int `json:"cache_size_mb"`
			MaxYears    int `json:"max_years"`
		} `json:"calendar_img"`
	} `json:"limits"`
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
