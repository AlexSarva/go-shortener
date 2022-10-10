package app

import (
	"AlexSarva/go-shortener/constant"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/storage"
	"AlexSarva/go-shortener/storage/storagefile"
	"AlexSarva/go-shortener/storage/storagelocal"
	"AlexSarva/go-shortener/storage/storagepg"
	"fmt"
	"log"
)

// Database interface for different types of databases
type Database struct {
	Repo storage.Repo
}

// NewStorage generate new instance of database
func NewStorage() *Database {

	cfg := constant.GlobalContainer.Get("server-config").(models.Config)

	if cfg.FileStorage == "" && cfg.Database == "" {
		Storage := storagelocal.NewURLLocalStorage()
		fmt.Println("Using In-Memory Database")
		return &Database{
			Repo: Storage,
		}
	} else if len(cfg.Database) > 0 {
		Storage := storagepg.NewPostgresDBConnection(cfg.Database)
		fmt.Println("Using PostgreSQL Database")
		return &Database{
			Repo: Storage,
		}
	} else {
		fmt.Println("Using file Database")
		Storage, err := storagefile.NewFileStorage(cfg.FileStorage)
		if err != nil {
			log.Fatal(err)
		}
		return &Database{
			Repo: Storage,
		}
	}
}
