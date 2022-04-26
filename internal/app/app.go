package app

import (
	"AlexSarva/go-shortener/storage"
	"AlexSarva/go-shortener/storage/storagelocal"
)

type Database struct {
	Repo storage.Repo
}

func NewDB() *Database {
	DBStorage := storagelocal.NewURLLocalStorage()
	return &Database{
		Repo: DBStorage,
	}
}
