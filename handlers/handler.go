package handlers

import (
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/utils"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const ShortLen int = 10

// Обработка сжатых запросов
func readBodyBytes(r *http.Request) (io.ReadCloser, error) {
	// GZIP decode
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		// Read body
		bodyBytes, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			return nil, readErr
		}
		defer r.Body.Close()

		log.Println("Получен Сжатый Body")

		newR, gzErr := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			log.Println(gzErr)
			return nil, gzErr
		}
		defer newR.Close()

		//bb, err2 := ioutil.ReadAll(r)
		//if err2 != nil {
		//	return nil, err2
		//}
		log.Println("Возвращен нормальный Body")
		return newR, nil
	} else {
		log.Println("Получен Обычный Body")
		// Not compressed
		return r.Body, nil
	}
}

// Дополнительный обработчик ошибок
func errorResponse(w http.ResponseWriter, message, errContentType string, httpStatusCode int) {
	w.Header().Set("Content-Type", errContentType)
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

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
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Add("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

		_, err := w.Write([]byte(longURL))
		if err != nil {
			log.Println("Something wrong", err)
		}

	}
}

func MakeShortURLHandler(cfg *models.Config, database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains("text/plain, text/xml, text/plain", r.Header.Get("Content-Type")) {
			errorResponse(w, "Content Type is not text", "text/plain", http.StatusUnsupportedMediaType)
			return
		}

		b, err := readBodyBytes(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
		}

		body, bodyErr := io.ReadAll(b)
		if bodyErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(bodyErr.Error())
		}

		rawURL := string(body)
		log.Println(rawURL)
		if utils.ValidateURL(rawURL) {
			id := utils.ShortURLGenerator(ShortLen)
			dbErr := database.Repo.InsertURL(id, rawURL, cfg.BaseURL)
			if dbErr != nil {
				log.Println(dbErr)
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			log.Println("URL write to DB")

			newShortURL, _ := database.Repo.GetURL(id)
			_, err := w.Write([]byte(newShortURL.ShortURL))
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

func MakeShortURLByJSON(cfg *models.Config, database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		headerContentTtype := r.Header.Get("Content-Type")

		if headerContentTtype != "application/json" {
			errorResponse(w, "Content Type is not application/json", "application/json", http.StatusUnsupportedMediaType)
			return
		}

		var newURL models.NewURL
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&newURL)

		if err != nil {
			if errors.As(err, &unmarshalErr) {
				errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				errorResponse(w, "Bad Request "+err.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		if utils.ValidateURL(newURL.URL) {
			id := utils.ShortURLGenerator(ShortLen)
			dbErr := database.Repo.InsertURL(id, newURL.URL, cfg.BaseURL)
			if dbErr != nil {
				log.Println(dbErr)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			log.Println("URL write to DB")

			newShortURL, _ := database.Repo.GetURL(id)

			resultURL := models.ResultURL{
				Result: newShortURL.ShortURL,
			}
			bodyURL, bodyErr := json.Marshal(resultURL)
			if bodyErr != nil {
				panic(bodyErr)
			}
			_, err := w.Write(bodyURL)
			if err != nil {
				log.Println("Something wrong", err)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("It's not URL!"))
			if err != nil {
				log.Println("Something wrong", err)
			}
		}
	}
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

var gzipContentTypes = "application/javascript, application/json, text/css, text/html, text/plain, text/xml"

func GzipHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(gzipContentTypes, r.Header.Get("Content-Type")) {
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
	return http.HandlerFunc(fn)
}

func MyHandler(cfg *models.Config, database *app.Database) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(GzipHandler)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain"))
	r.Use(middleware.Compress(5, gzipContentTypes))

	r.Get("/{id}", GetRedirectURL(database))
	r.Post("/", MakeShortURLHandler(cfg, database))
	r.Post("/api/shorten", MakeShortURLByJSON(cfg, database))

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
