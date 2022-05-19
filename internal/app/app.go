package app

import (
	"AlexSarva/go-shortener/storage"
	"AlexSarva/go-shortener/storage/storagefile"
	"AlexSarva/go-shortener/storage/storagelocal"
	"fmt"
	"log"
)

type Database struct {
	Repo storage.Repo
}

func NewDB(filePath string) *Database {
	if filePath == "" {
		Storage := storagelocal.NewURLLocalStorage()
		fmt.Println("Using In-Memory Database")
		return &Database{
			Repo: Storage,
		}
	} else {
		fmt.Println("Using file Database")
		Storage, err := storagefile.NewFileStorage(filePath)
		if err != nil {
			log.Fatal(err)
		}
		return &Database{
			Repo: Storage,
		}
	}
}
