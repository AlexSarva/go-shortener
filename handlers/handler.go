package handlers

import (
	"AlexSarva/go-shortener/constant"
	"AlexSarva/go-shortener/crypto"
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/storage"
	"AlexSarva/go-shortener/utils"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ShortLen length of the short URL
const ShortLen int = 10

// GzipContentTypes accepted content types
const GzipContentTypes = "application/x-gzip, application/javascript, application/json, text/css, text/html, text/plain, text/xml"

// ErrNotValidCookie expects when cookie not found or not valid
var ErrNotValidCookie = errors.New("valid cookie does not found")

// ReadBodyBytes handling compressed requests
func ReadBodyBytes(r *http.Request) (io.ReadCloser, error) {
	// GZIP decode
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		// Read body
		bodyBytes, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			return nil, readErr
		}
		defer r.Body.Close()

		log.Println("compressed request")

		newR, gzErr := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			log.Println(gzErr)
			return nil, gzErr
		}
		defer newR.Close()

		return newR, nil
	} else {
		log.Println("no compressed request")
		return r.Body, nil
	}
}

// ErrorResponse Additional error handler
func ErrorResponse(w http.ResponseWriter, message, errContentType string, httpStatusCode int) {
	w.Header().Set("Content-Type", errContentType)
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

// PingDB check connection to DB
//
// Possible response codes:
// 200 — connection OK;
// 500 is an internal server error.
func PingDB(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ping := database.Repo.Ping()
		if ping {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

	}
}

// GetRedirectURL accepts a shortened URL identifier as a URL parameter and returns a response with a 307 code
// and the original URL in the Location HTTP header
//
// Handler: GET /{id}
//
// Possible response codes:
// 201 — successful add links;
// 400 - wrong link format;
// 410 - link deleted from DB;
// 500 is an internal server error.
func GetRedirectURL(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log.Println(id)
		res, er := database.Repo.GetURL(id)
		if er != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("No such shortlink!"))
			if err != nil {
				log.Println("Something wrong", err)
			}
			return
		}
		if res.Deleted == 1 {
			w.WriteHeader(http.StatusGone)
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

// GetUserURLs return to the user all the URLs they have ever shortened
//
// Handler: GET /api/user/urls
//
// Respond format:
//
//		[
//	   {
//	       "short_url": "http://localhost:9001/EpRZQwytfH",
//	       "original_url": "https://t.me/moscowach"
//	   },
//	   {
//	       "short_url": "http://localhost:9001/DnZneSjGhd",
//	       "original_url": "https://clickhouse.com"
//	   },
//		]
//
// Possible response codes:
// 200 - status OK;
// 204 - no content by cookie;
func GetUserURLs(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID, userIDErr := GetCookie(r)
		if userIDErr != nil {
			log.Println(userIDErr)
			w.WriteHeader(http.StatusNoContent)
			return

		}
		res, er := database.Repo.GetUserURLs(userID.String())
		if er != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resultList, resultListErr := json.Marshal(res)
		if resultListErr != nil {
			panic(resultListErr)
		}
		_, err := w.Write(resultList)
		if err != nil {
			log.Println("Something wrong", err)
		}
		w.WriteHeader(http.StatusOK)
	}

}

// MakeShortURLHandler accepts URL string to shortener in the request body and returns a response with a 201 code
// and shortened URL as a text string in the body.
//
// Handler: POST /
//
// Accept Content-Type: text/plain, text/xml, text/plain, text/plain; charset=utf-8, application/x-gzip
//
// Possible response codes:
// 201 — successful add link;
// 400 - wrong link format;
// 409 - link already exists in DB;
// 415 - wrong request type;
// 500 is an internal server error.
func MakeShortURLHandler(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains("text/plain, text/xml, text/plain, text/plain; charset=utf-8, application/x-gzip", r.Header.Get("Content-Type")) {
			ErrorResponse(w, "Content Type is not a text/plain or application/x-gzip", "text/plain", http.StatusUnsupportedMediaType)
			return
		}

		cfg := constant.GlobalContainer.Get("server-config").(models.Config)

		cookie, _ := r.Cookie("session")
		userID, _ := crypto.Decrypt(cookie.Value, crypto.SecretKey)

		b, err := ReadBodyBytes(r)
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
			shortURL := utils.CreateShortUrl(cfg.BaseURL, id)
			dbErr := database.Repo.InsertURL(id, rawURL, shortURL, userID.String())
			if dbErr == storage.ErrDuplicatePK {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				existShortURL, _ := database.Repo.GetURLByRaw(rawURL)
				_, err := w.Write([]byte(existShortURL.ShortURL))
				if err != nil {
					log.Println("Something wrong", err)
				}
				return
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

// MakeShortURLByJSON accepts a JSON object  in the request body and returns an object {"result":"<shorten_url>"} in response.
//
// Handler: POST /api/shorten
//
// Request format:
//
//	{"url":"yandex.ru"}
//
// Respond format:
//
//	{"result":"http://localhost:9001/QBRfxXmcYp"}
//
// Possible response codes:
// 201 — successful add links;
// 400 - wrong link format;
// 409 - link already exists in DB;
// 415 - wrong request type;
// 500 is an internal server error.
func MakeShortURLByJSON(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		headerContentTtype := r.Header.Get("Content-Type")

		if !strings.Contains("application/json, application/x-gzip", headerContentTtype) {
			ErrorResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusUnsupportedMediaType)
			return
		}

		cfg := constant.GlobalContainer.Get("server-config").(models.Config)

		cookie, _ := r.Cookie("session")
		userID, _ := crypto.Decrypt(cookie.Value, crypto.SecretKey)

		var newURL models.NewURL
		var unmarshalErr *json.UnmarshalTypeError

		b, err := ReadBodyBytes(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&newURL)

		if err != nil {
			if errors.As(err, &unmarshalErr) {
				ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				ErrorResponse(w, "Bad Request "+err.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		if utils.ValidateURL(newURL.URL) {
			id := utils.ShortURLGenerator(ShortLen)
			shortURL := utils.CreateShortUrl(cfg.BaseURL, id)
			dbErr := database.Repo.InsertURL(id, newURL.URL, shortURL, userID.String())
			if dbErr == storage.ErrDuplicatePK {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				existShortURL, _ := database.Repo.GetURLByRaw(newURL.URL)

				existURL := models.ResultURL{
					Result: existShortURL.ShortURL,
				}
				bodyURL, bodyErr := json.Marshal(existURL)
				if bodyErr != nil {
					panic(bodyErr)
				}
				_, err := w.Write(bodyURL)
				if err != nil {
					log.Println("Something wrong", err)
				}
				return
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

// MakeBatchURLByJSON accept in the request body a set of URLs to shorten in the format
//
// Handler: POST /api/shorten/batch
//
// Request format:
//
//	[{
//	  "correlation_id": "1",
//	  "original_url": "yandex.ru"
//	},
//	...]
//
// Respond format:
//
//	[{
//		"correlation_id": "1",
//		"short_url": "http://localhost:9001/QBRfxXmcYp"
//	},
//	...]
//
// Possible response codes:
// 201 — successful add links;
// 400 - wrong link format;
// 415 - wrong request type;
// 500 is an internal server error.
func MakeBatchURLByJSON(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		headerContentType := r.Header.Get("Content-Type")

		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			ErrorResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusUnsupportedMediaType)
			return
		}

		cfg := constant.GlobalContainer.Get("server-config").(models.Config)

		var rawBatchURL []models.RawBatchURL
		var resultBatchURL []models.ResultBatchURL
		var insertBatchURL []models.URL
		var unmarshalErr *json.UnmarshalTypeError

		b, err := ReadBodyBytes(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&rawBatchURL)

		if err != nil {
			log.Println(err)
			if errors.As(err, &unmarshalErr) {
				ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				ErrorResponse(w, "Bad Request "+err.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		cookie, _ := r.Cookie("session")
		userID, _ := crypto.Decrypt(cookie.Value, crypto.SecretKey)

		for _, urlInfo := range rawBatchURL {
			if utils.ValidateURL(urlInfo.RawURL) {
				id := utils.ShortURLGenerator(ShortLen)
				shortURL := utils.CreateShortUrl(cfg.BaseURL, id)
				currentURLInsert := models.URL{
					ID:       id,
					RawURL:   urlInfo.RawURL,
					ShortURL: shortURL,
					Created:  time.Now(),
					UserID:   userID.String(),
				}
				currentURLResult := models.ResultBatchURL{
					CorrelationID: urlInfo.CorrelationID,
					ShortURL:      shortURL,
				}
				insertBatchURL = append(insertBatchURL, currentURLInsert)
				resultBatchURL = append(resultBatchURL, currentURLResult)

			} else {
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte("It's not URL!"))
				if err != nil {
					log.Println("Something wrong", err)
				}
			}
		}
		dbErr := database.Repo.InsertMany(insertBatchURL)
		if dbErr != nil {
			log.Println(dbErr)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		bodyURL, bodyErr := json.Marshal(resultBatchURL)
		if bodyErr != nil {
			panic(bodyErr)
		}
		_, writeErr := w.Write(bodyURL)
		if writeErr != nil {
			log.Println("Something wrong", err)
		}
	}
}

// AddDeleteURLs async delete url from DB using channels
func AddDeleteURLs(urls models.DeleteURL, deleteCh chan models.DeleteURL) {
	deleteCh <- urls
}

// DeleteAsync accept in the request body a set of URLs to shorten in the format
//
// Handler: DELETE /api/user/urls
//
// Request format:
// [ "a", "b", "c", "d", ...]
//
// Possible response codes:
// 202 — accepted;
// 400 - wrong link format;
// 415 - wrong request type;
func DeleteAsync(deleteCh chan models.DeleteURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		headerContentTtype := r.Header.Get("Content-Type")

		if !strings.Contains("application/json, application/x-gzip", headerContentTtype) {
			ErrorResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusUnsupportedMediaType)
			return
		}

		var deleteBatchURL []string
		var unmarshalErr *json.UnmarshalTypeError
		b, err := ReadBodyBytes(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&deleteBatchURL)

		if err != nil {
			log.Println(err)
			if errors.As(err, &unmarshalErr) {
				ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				ErrorResponse(w, "Bad Request "+err.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		cookie, _ := r.Cookie("session")
		userID, _ := crypto.Decrypt(cookie.Value, crypto.SecretKey)

		go AddDeleteURLs(models.DeleteURL{
			UserID: userID.String(),
			URLs:   deleteBatchURL,
		}, deleteCh)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		log.Printf("%+v\n", deleteBatchURL)
	}
}

// gzipWriter struct to write response in compressed form
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write response in compressed form
func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// GzipHandler middleware handle g-zipped requests
func GzipHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(GzipContentTypes, r.Header.Get("Content-Type")) {
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

// MyHandler base handler for all routes and middleware
func MyHandler(database *app.Database, deleteCh chan models.DeleteURL) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CookieHandler)
	r.Use(GzipHandler)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, GzipContentTypes))
	r.Mount("/debug", middleware.Profiler())
	r.Get("/ping", PingDB(database))
	r.Get("/{id}", GetRedirectURL(database))
	r.Post("/", MakeShortURLHandler(database))
	r.Post("/api/shorten", MakeShortURLByJSON(database))
	r.Post("/api/shorten/batch", MakeBatchURLByJSON(database))
	r.Get("/api/user/urls", GetUserURLs(database))
	r.Delete("/api/user/urls", DeleteAsync(deleteCh))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, nfErr := w.Write([]byte("route does not exist"))
		if nfErr != nil {
			log.Println(nfErr)
		}
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, naErr := w.Write([]byte("sorry, only GET, POST and DELETE methods are supported."))
		if naErr != nil {
			log.Println(naErr)
		}
	})
	return r
}
