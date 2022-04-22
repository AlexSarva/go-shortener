package main

import (
	"AlexSarva/go-shortener/server"
	"AlexSarva/go-shortener/storage"
	"log"
)

func main() {
	DB := *storage.InitDB()
	MainApp := server.NewMyServer(8080, DB)
	if err := MainApp.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
