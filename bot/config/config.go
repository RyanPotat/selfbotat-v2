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

type TwitchConfig struct {
	Login string `json:"login"`
	Id string `json:"id"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
	TLAHeaders map[string]string `json:"tlaHeaders"`
	TLAURI string `json:"tlaUri"`
}

type BotConfig struct {
	Prefix string `json:"prefix"`
	Postgres PostgresConfig `json:"postgres"`
	Twitch TwitchConfig `json:"twitch"`
}

var Config BotConfig

func LoadConfig() (BotConfig, error) {
	fillData, err := os.ReadFile("config.json")
	if err != nil {
		return BotConfig{}, fmt.Errorf("error loading config: %w", err)
	}

	err = json.Unmarshal(fillData, &Config)
	if err != nil {
		return BotConfig{}, fmt.Errorf("error parsing config: %w", err)
	}

	return Config, nil
}