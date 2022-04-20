package main

import (
	"go-shortener/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.BodyHandler)
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
