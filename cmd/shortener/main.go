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
	cfg          models.Config
	JSONConfig   models.JSONConfig
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

func init() {
	flag.StringVar(&cfg.ServerAddress, "a", "", "host:port to listen on")
	flag.StringVar(&cfg.BaseURL, "b", "", "base host:port for short link")
	flag.StringVar(&cfg.FileStorage, "f", "", "filepath of short links file storage")
	flag.StringVar(&cfg.Database, "d", "", "database config")
	flag.StringVar(&JSONConfig.DSN, "c", "", "JSON config")
	flag.BoolVar(&cfg.EnableHTTPS, "s", false, "enable HTTPS")
}

func main() {
	version()

	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ServerAddress: %v, BaseURL: %v, FileStorage: %v, EnableHTTPS: %v", cfg.ServerAddress, cfg.BaseURL, cfg.FileStorage, cfg.EnableHTTPS)

	// Перезаписываем из параметров запуска
	flag.Parse()

	if cfg == (models.Config{}) {
		if configFilename := JSONConfig.DSN; configFilename != "" {
			JSONErr := models.ReadJSONConfig(&cfg, configFilename)
			if JSONErr != nil {
				log.Println(JSONErr)
			}
		}
		cfg.BaseURL = "http://localhost:8080"
		cfg.ServerAddress = "localhost:8080"
	}

	log.Printf("ServerAddress: %v, BaseURL: %v, FileStorage: %v, EnableHTTPS: %v", cfg.ServerAddress, cfg.BaseURL, cfg.FileStorage, cfg.EnableHTTPS)

	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}

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
