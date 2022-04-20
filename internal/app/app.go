package app

import "go-shortener/storage"

func InitDB() *storage.UrlLocalStorage {
	Db := storage.NewUrlLocalStorage()
	return Db
}
