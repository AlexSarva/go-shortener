package handlers

import (
	"AlexSarva/go-shortener/checker"
	"AlexSarva/go-shortener/generator"
	"AlexSarva/go-shortener/storage"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

const ShortLen int = 5

func MyGetHandle(database *storage.UrlLocalStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		res, er := database.Get(id)
		if er != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("No such shortlink!"))
			if err != nil {
				log.Println("Something wrong", err)
			}
		} else {
			longUrl := res.RawUrl
			w.Header().Set("content-type", "text/plain")
			w.Header().Add("Location", longUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)

			_, err := w.Write([]byte(longUrl))
			if err != nil {
				log.Println("Something wrong", err)
			}
		}
	}
}

func MyPostHandle(database *storage.UrlLocalStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		fmt.Println(b)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
		}
		rawUrl := string(b)
		if checker.CheckUrl(rawUrl) {
			shortUrl := generator.ShortUrlGenerator(ShortLen)
			dbErr := database.Insert(rawUrl, shortUrl)
			if dbErr != nil {
				log.Println(dbErr)
			}
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			log.Println("URL write to DB")
			_, err := w.Write([]byte("http://localhost:8080/" + shortUrl))
			if err != nil {
				log.Println("Something wrong", err)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("It's not Url!"))
			if err != nil {
				log.Println("Something wrong", err)
			}
			log.Println("It's not Url!")
		}
	}
}
