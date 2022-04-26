package handlers

import (
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/utils"

	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const ShortLen int = 5

func GetRedirectURL(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		res, er := database.Repo.GetURL(id)
		if er != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("No such shortlink!"))
			if err != nil {
				log.Println("Something wrong", err)
			}
			return
		}
		longURL := res.RawURL
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.Header().Add("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

		_, err := w.Write([]byte(longURL))
		if err != nil {
			log.Println("Something wrong", err)
		}

	}
}

func MakeShortURLHandler(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		fmt.Println(b)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
		}
		rawURL := string(b)
		if utils.ValidateURL(rawURL) {
			shortURL := utils.ShortURLGenerator(ShortLen)
			dbErr := database.Repo.InsertURL(rawURL, shortURL)
			if dbErr != nil {
				log.Println(dbErr)
			}
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			log.Println("URL write to DB")
			_, err := w.Write([]byte("http://localhost:8080/" + shortURL))
			if err != nil {
				log.Println("Something wrong", err)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("It's not URL!"))
			if err != nil {
				log.Println("Something wrong", err)
			}
			log.Println("It's not URL!")
		}
	}
}

func MyHandler(database *app.Database) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", GetRedirectURL(database))
	r.Post("/", MakeShortURLHandler(database))

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
	return r
}
