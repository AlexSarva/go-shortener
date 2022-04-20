package handlers

import (
	"go-shortener/checker"
	"go-shortener/generator"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type MyUrl struct {
	ID      int
	LongUrl string
	Created time.Time
}

var Urls = make(map[string]MyUrl)
var urlId int = 1

const ShortLen int = 5

func BodyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
		}
		rawUrl := string(b)
		if checker.CheckUrl(rawUrl) {
			shortUrl := generator.ShortUrlGenerator(ShortLen)
			var UrlData MyUrl
			UrlData.ID = urlId
			UrlData.LongUrl = rawUrl
			UrlData.Created = time.Now()
			Urls[shortUrl] = UrlData
			urlId += 1

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

	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/")
		res, er := Urls[id]
		if er == false {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("No such shortlink!"))
			if err != nil {
				log.Println("Something wrong", err)
			}
		}

		longUrl := res.LongUrl
		w.Header().Set("content-type", "text/plain")
		w.Header().Add("Location", longUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

		_, err := w.Write([]byte(longUrl))
		if err != nil {
			log.Println("Something wrong", err)
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Sorry, only GET and POST methods are supported."))
		if err != nil {
			log.Println("Something wrong", err)
		}
	}

}
