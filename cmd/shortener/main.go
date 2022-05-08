package main

import (
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/server"
	"github.com/caarlos0/env/v6"

	"log"
)

func main() {
	var cfg models.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	DB := *app.NewDB()
	MainApp := server.NewMyServer(&cfg, &DB)
	if err := MainApp.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
