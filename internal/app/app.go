package app

import (
	"AlexSarva/go-shortener/handlers"
	"AlexSarva/go-shortener/storage"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func InitDB() *storage.UrlLocalStorage {
	db := *storage.NewUrlLocalStorage()
	return &db
}

type App struct {
	httpServer *http.Server
	router     *chi.Mux
}

func NewApp(port int) *App {
	database := *InitDB()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", handlers.MyGetHandle(&database))
	r.Post("/", handlers.MyPostHandle(&database))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, nfErr := w.Write([]byte("route does not exist"))
		if nfErr != nil {
			log.Println(nfErr)
		}
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, naErr := w.Write([]byte("sorry, only GET and POST methods are supported."))
		if naErr != nil {
			log.Println(naErr)
		}
	})

	server := http.Server{
		Addr:    "localhost:" + strconv.Itoa(port),
		Handler: r,
	}
	return &App{
		httpServer: &server,
	}
}

func (a *App) Run() error {
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
