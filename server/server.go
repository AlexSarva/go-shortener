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
	"syscall"
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

	idleConnsClosed := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-quit
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := a.httpServer.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown: %v", err)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	if cfg.EnableHTTPS {
		log.Printf("Web-server started at https://%s", addr)
		go func() {
			if err := a.httpServer.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"); err != http.ErrServerClosed {
				log.Fatalf("Failed to listen and serve TLS: %+v", err)
			}
		}()
	} else {
		log.Printf("Web-server started at http://%s", addr)
		go func() {
			if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("Failed to listen and serve: %+v", err)
			}
		}()
	}

	<-idleConnsClosed
	//
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	log.Println("Server Shutdown gracefully")
	return a.httpServer.Shutdown(ctx)
}
