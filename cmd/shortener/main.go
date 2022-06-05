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
	log.Printf("ServerAddress: %v, BaseURL: %v, FileStorage: %v", cfg.ServerAddress, cfg.BaseURL, cfg.FileStorage)
	// Перезаписываем из параметров запуска
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "host:port to listen on")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base host:port for short link")
	flag.StringVar(&cfg.FileStorage, "f", cfg.FileStorage, "filepath of short links file storage")
	flag.StringVar(&cfg.Database, "d", cfg.Database, "database config")
	flag.Parse()
	log.Printf("ServerAddress: %v, BaseURL: %v, FileStorage: %v", cfg.ServerAddress, cfg.BaseURL, cfg.FileStorage)
	DB := *app.NewStorage(cfg.FileStorage, cfg.Database)
	ping := DB.Repo.Ping()
	log.Println(ping)
	MainApp := server.NewMyServer(&cfg, &DB)
	if err := MainApp.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
