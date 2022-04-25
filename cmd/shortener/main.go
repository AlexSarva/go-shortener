package main

import (
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/server"
	"log"
)

func main() {
	DB := app.NewDB()
	MainApp := server.NewMyServer(8080, DB)
	if err := MainApp.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
