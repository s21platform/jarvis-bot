package model

import (
	"encoding/json"
	"log"
)

type Config struct {
	Value json.RawMessage `db:"value"`
}

type BirthdayPoolConfig struct {
	Enabled                bool     `json:"enabled"`
	HourStart              int      `json:"hour_start"`
	HourEnd                int      `json:"hour_finish"`
	WhiteList              []string `json:"white_list"`
	PendingIntervalMinutes int      `json:"pending_interval_minutes"`
}

func (c *Config) GetBirthdayPoolConfig() BirthdayPoolConfig {
	var config BirthdayPoolConfig
	err := json.Unmarshal(c.Value, &config)
	if err != nil {
		log.Println("failed to unmarshal config, return default config: ", err)
		return BirthdayPoolConfig{
			Enabled:                false,
			HourStart:              10,
			HourEnd:                20,
			WhiteList:              []string{"garroshm"},
			PendingIntervalMinutes: 10,
		}
	}
	return config
}
