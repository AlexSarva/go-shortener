package main

import (
	"AlexSarva/go-shortener/constant"
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/server"
	"flag"
	"log"
	_ "net/http/pprof"

	"github.com/caarlos0/env/v6"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func version() {
	log.Printf("Build version: %s\n", buildVersion)
	log.Printf("Build date: %s\n", buildDate)
	log.Printf("Build commit: %s\n", buildCommit)
}

func readChan(delCh chan models.DeleteURL, database *app.Database) {
	for v := range delCh {
		log.Println(v)
		err := database.Repo.Delete(v.UserID, v.URLs)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	version()
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

	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}

	log.Printf("ServerAddress: %v, BaseURL: %v, FileStorage: %v", cfg.ServerAddress, cfg.BaseURL, cfg.FileStorage)
	DB := *app.NewStorage()
	ping := DB.Repo.Ping()
	log.Println(ping)
	deleteCh := make(chan models.DeleteURL)
	MainApp := server.NewMyServer(&DB, deleteCh)
	go readChan(deleteCh, &DB)
	if err := MainApp.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
