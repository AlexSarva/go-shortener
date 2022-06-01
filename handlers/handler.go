package handlers

import (
	"AlexSarva/go-shortener/crypto"
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/storage"
	"AlexSarva/go-shortener/utils"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const ShortLen int = 10

var NotValidCookieErr = errors.New("valid cookie does not found")

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

		log.Println("Получен сжатый запрос")

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
		//log.Println("Возвращен нормальный Body")
		return newR, nil
	} else {
		log.Println("Получен несжатый запрос")
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

func GetUserURLs(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, userIDErr := getCookie(r)
		if userIDErr != nil {
			log.Println(userIDErr)
			w.WriteHeader(http.StatusNoContent)
			return

		}
		res, er := database.Repo.GetUserURLs(userID.String())
		//log.Printf("%+v\n", res)
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

}

func MakeShortURLHandler(cfg *models.Config, database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains("text/plain, text/xml, text/plain, text/plain; charset=utf-8, application/x-gzip", r.Header.Get("Content-Type")) {
			errorResponse(w, "Content Type is not a text/plain or application/x-gzip", "text/plain", http.StatusUnsupportedMediaType)
			return
		}
		cookie, _ := r.Cookie("session")
		userID, _ := crypto.Decrypt(cookie.Value, []byte("test"))

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
			dbErr := database.Repo.InsertURL(id, rawURL, cfg.BaseURL, userID.String())
			if dbErr == storage.DuplicatePKErr {
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

func MakeShortURLByJSON(cfg *models.Config, database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		headerContentTtype := r.Header.Get("Content-Type")

		if !strings.Contains("application/json, application/x-gzip", headerContentTtype) {
			errorResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusUnsupportedMediaType)
			return
		}

		cookie, _ := r.Cookie("session")
		userID, _ := crypto.Decrypt(cookie.Value, []byte("test"))

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
			dbErr := database.Repo.InsertURL(id, newURL.URL, cfg.BaseURL, userID.String())
			if dbErr == storage.DuplicatePKErr {
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

func MakeBatchURLByJSON(cfg *models.Config, database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		headerContentTtype := r.Header.Get("Content-Type")

		if !strings.Contains("application/json, application/x-gzip", headerContentTtype) {
			errorResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusUnsupportedMediaType)
			return
		}

		var rawBatchURL []models.RawBatchURL
		var resultBatchURL []models.ResultBatchURL
		var insertBatchURL []models.URL
		//var newURL models.NewURL
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
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
				errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				errorResponse(w, "Bad Request "+err.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		cookie, _ := r.Cookie("session")
		userID, _ := crypto.Decrypt(cookie.Value, []byte("test"))

		for _, urlInfo := range rawBatchURL {
			if utils.ValidateURL(urlInfo.RawURL) {
				id := utils.ShortURLGenerator(ShortLen)
				shortURL := cfg.BaseURL + "/" + id
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

		log.Printf("%+v\n", insertBatchURL)
		log.Printf("%+v\n", resultBatchURL)
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

var gzipContentTypes = "application/x-gzip, application/javascript, application/json, text/css, text/html, text/plain, text/xml"

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

func GenerateCookie(userID uuid.UUID) http.Cookie {
	session := crypto.Encrypt(userID, []byte("test"))
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "session", Value: session, Expires: expiration, Path: "/"}
	return cookie
}

func getCookie(r *http.Request) (uuid.UUID, error) {
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		log.Println(cookieErr)
		return uuid.UUID{}, NotValidCookieErr
	}
	userID, cookieDecryptErr := crypto.Decrypt(cookie.Value, []byte("test"))
	if cookieDecryptErr != nil {
		return uuid.UUID{}, cookieDecryptErr
	}
	return userID, nil

}

func CookieHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, userIDErr := getCookie(r)
		if userIDErr != nil {
			log.Println(userIDErr)
			userCookie := GenerateCookie(uuid.New())
			log.Println(userCookie)
			r.AddCookie(&userCookie)
			http.SetCookie(w, &userCookie)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func MyHandler(cfg *models.Config, database *app.Database) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CookieHandler)
	r.Use(GzipHandler)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, gzipContentTypes))

	r.Get("/ping", PingDB(database))
	r.Get("/{id}", GetRedirectURL(database))
	r.Post("/", MakeShortURLHandler(cfg, database))
	r.Post("/api/shorten", MakeShortURLByJSON(cfg, database))
	r.Post("/api/shorten/batch", MakeBatchURLByJSON(cfg, database))
	r.Get("/api/user/urls", GetUserURLs(database))

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
