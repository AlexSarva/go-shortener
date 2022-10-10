package server

import (
	"AlexSarva/go-shortener/constant"
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

// MyServer implementation of custom server
type MyServer struct {
	httpServer *http.Server
}

// NewMyServer Initializing new server instance
func NewMyServer(database *app.Database, deleteCh chan models.DeleteURL) *MyServer {

	cfg := constant.GlobalContainer.Get("server-config").(models.Config)

	handler := handlers.MyHandler(database, deleteCh)
	server := http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
	}
	return &MyServer{
		httpServer: &server,
	}
}

// Run method that starts the server
func (a *MyServer) Run() error {
	addr := a.httpServer.Addr

	cfg := constant.GlobalContainer.Get("server-config").(models.Config)

	if cfg.EnableHTTPS {
		log.Printf("Web-server started at https://%s", addr)
		go func() {
			if err := a.httpServer.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"); err != nil {
				log.Fatalf("Failed to listen and serve TLS: %+v", err)
			}
		}()
	} else {
		log.Printf("Web-server started at http://%s", addr)
		go func() {
			if err := a.httpServer.ListenAndServe(); err != nil {
				log.Fatalf("Failed to listen and serve: %+v", err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
