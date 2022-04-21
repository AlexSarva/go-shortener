package main

import (
	"AlexSarva/go-shortener/internal/app"
	"log"
)

func main() {
	MainApp := app.NewApp()
	if err := MainApp.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
