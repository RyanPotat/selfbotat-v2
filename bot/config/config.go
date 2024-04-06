package config

import (
	"encoding/json"
	"fmt"
	"os"

	"selfbotat-v2/bot/types"
)

var Config types.BotConfig

func LoadConfig() error {
	fillData, err := os.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	err = json.Unmarshal(fillData, &Config)
	if err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	return nil
}