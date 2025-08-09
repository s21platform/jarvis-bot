package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s21platform/jarvis-bot/internal/api"
	"github.com/s21platform/jarvis-bot/internal/config"
	api2 "github.com/s21platform/jarvis-bot/internal/generated"
	"github.com/s21platform/jarvis-bot/internal/repository/postgres"
)

func main() {
	cfg := config.MustLoadConfig()
	log.Println(cfg.Service.Url)

	repo := postgres.New(cfg)
	h := api.New(cfg, repo, cfg.Service.Url)

	r := chi.NewRouter()

	api2.HandlerFromMux(h, r)

	log.Println("Starting server", cfg.Service.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Service.Port), r))
}
