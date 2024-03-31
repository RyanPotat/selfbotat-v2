package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Config struct {
	Login string `json:"login"`
	Id string `json:"id"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
	Prefix string `json:"prefix"`
	Postgres PostgresConfig `json:"postgres"`
}

var cfg Config

func LoadConfig() (Config, error) {
	fillData, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, fmt.Errorf("error loading config: %w", err)
	}

	err = json.Unmarshal(fillData, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config: %w", err)
	}

	return cfg, nil
}

func GetConfig() Config {
	return cfg
}