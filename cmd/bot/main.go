package main

import (
	"github.com/s21platform/jarvis-bot/internal/config"
)

func main() {
	_ = config.MustLoadConfig()
}
