package server

import (
	"AlexSarva/go-shortener/handlers"
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type MyServer struct {
	httpServer *http.Server
}

func NewMyServer(cfg *models.Config, database *app.Database) *MyServer {
	handler := handlers.MyHandler(cfg, database)
	server := http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
	}
	return &MyServer{
		httpServer: &server,
	}
}

func (a *MyServer) Run() error {
	addr := a.httpServer.Addr
	log.Printf("Запуск веб-сервера на http://%s", addr)
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
