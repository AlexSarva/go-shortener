package main

import (
	"AlexSarva/go-shortener/internal/app"
	"log"
)

func main() {
	MainApp := app.NewApp(8080)
	if err := MainApp.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
