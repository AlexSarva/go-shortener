package app

import (
	"AlexSarva/go-shortener/storage"
	"AlexSarva/go-shortener/storage/storageLocal"
)

type Database struct {
	Repo storage.Repo
}

func NewDB() *Database {
	DBStorage := storageLocal.NewURLLocalStorage()
	return &Database{
		Repo: DBStorage,
	}
}
