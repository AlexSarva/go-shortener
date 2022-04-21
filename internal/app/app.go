package app

import (
	"AlexSarva/go-shortener/handlers"
	"AlexSarva/go-shortener/storage"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func InitDB() *storage.UrlLocalStorage {
	db := storage.NewUrlLocalStorage()
	return db
}

type App struct {
	Database   *storage.UrlLocalStorage
	httpServer *http.Server
	ServeMux   *http.ServeMux
}

func NewApp() *App {
	db := *InitDB()
	mux := *http.NewServeMux()
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: &mux,
	}
	return &App{
		Database:   &db,
		ServeMux:   &mux,
		httpServer: &server,
	}
}

func (a *App) Run() error {
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	a.ServeMux.Handle("/", handlers.MakeHandler(a.Database))
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
