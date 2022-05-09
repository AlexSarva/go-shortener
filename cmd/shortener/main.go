package main

import (
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/server"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

func main() {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)
	// Перезаписываем из параметров запуска
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "host:port to listen on")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base host:port for short link")
	flag.StringVar(&cfg.FileStorage, "f", cfg.FileStorage, "filepath of short links file storage")
	flag.Parse()
	log.Printf("%+v\n", cfg)
	DB := *app.NewDB(cfg.FileStorage)
	MainApp := server.NewMyServer(&cfg, &DB)
	if err := MainApp.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
