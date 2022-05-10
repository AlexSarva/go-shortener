package handlers_test

import (
	"AlexSarva/go-shortener/handlers"
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/utils"
	"github.com/caarlos0/env/v6"

	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMytHandler(t *testing.T) {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	database := app.NewDB(cfg.FileStorage)
	insErr := database.Repo.InsertURL("Hasfe", "https://codepen.io", cfg.BaseURL)
	log.Printf("%+v\n", database.Repo)
	if insErr != nil {
		log.Println(insErr)
	}
	type want struct {
		code            int
		location        string
		contentType     string
		contentEncoding string
		responseFormat  bool
		response        string
	}

	tests := []struct {
		name                   string
		request                string
		requestPath            string
		requestMethod          string
		requestBody            string
		requestCompressBody    []byte
		requestContentType     string
		requestAcceptEncoding  string
		requestContentEncoding string
		want                   want
	}{
		{
			name:          fmt.Sprintf("%s positive test #1", http.MethodGet),
			request:       "Hasfe",
			requestMethod: http.MethodGet,
			requestPath:   "/",
			want: want{
				code:        http.StatusTemporaryRedirect,
				location:    "https://codepen.io",
				contentType: "text/plain",
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodGet),
			request:       "",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #2", http.MethodGet),
			request:       "sometext",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s positive test #1", http.MethodPost),
			requestBody:        "https://codepen.io",
			requestPath:        "/",
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			want: want{
				code:           http.StatusCreated,
				responseFormat: true,
				contentType:    "text/plain",
			},
		},
		{
			name:               fmt.Sprintf("%s positive test #2", http.MethodPost),
			requestBody:        "www.google.com",
			requestPath:        "/",
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			want: want{
				code:           http.StatusCreated,
				responseFormat: true,
				contentType:    "text/plain",
			},
		},
		{
			name:               fmt.Sprintf("%s positive test #3", http.MethodPost),
			requestBody:        "www.url-with-querystring.com/?url=has-querystring",
			requestPath:        "/",
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			want: want{
				code:           http.StatusCreated,
				responseFormat: true,
				contentType:    "text/plain",
			},
		},
		{
			name:               fmt.Sprintf("%s positive test #4", http.MethodPost),
			requestBody:        "www.url-with-querystring.com/?url=has-querystring",
			requestPath:        "/",
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			want: want{
				code:        http.StatusUnsupportedMediaType,
				contentType: "text/plain",
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodPost),
			requestBody:        "https://codepen",
			requestPath:        "/",
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #2", http.MethodPost),
			requestBody:        "something",
			requestPath:        "/",
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #3", http.MethodPost),
			requestBody:        "",
			requestPath:        "/",
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodPut),
			requestBody:        "something",
			requestPath:        "/",
			request:            "something",
			requestContentType: "text/plain",
			requestMethod:      http.MethodPut,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodConnect),
			requestBody:        "something",
			requestPath:        "/",
			request:            "something",
			requestContentType: "text/plain",
			requestMethod:      http.MethodConnect,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodDelete),
			requestBody:        "something",
			requestPath:        "/",
			request:            "something",
			requestContentType: "text/plain",
			requestMethod:      http.MethodDelete,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodHead),
			requestBody:        "something",
			requestPath:        "/",
			request:            "something",
			requestContentType: "text/plain",
			requestMethod:      http.MethodHead,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodOptions),
			requestBody:        "something",
			requestPath:        "/",
			request:            "something",
			requestContentType: "text/plain",
			requestMethod:      http.MethodOptions,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodPatch),
			requestBody:        "something",
			requestPath:        "/",
			request:            "something",
			requestContentType: "text/plain",
			requestMethod:      http.MethodPatch,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodTrace),
			requestBody:        "something",
			requestPath:        "/",
			request:            "something",
			requestContentType: "text/plain",
			requestMethod:      http.MethodTrace,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", "CUSTOM"),
			requestBody:        "something",
			requestPath:        "/",
			request:            "something",
			requestMethod:      "CUSTOM",
			requestContentType: "text/plain",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s JSON positive test #1", http.MethodPost),
			requestBody:        `{"url":"https://codepen.io"}`,
			requestPath:        "/api/shorten",
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			want: want{
				code:           http.StatusCreated,
				responseFormat: true,
				contentType:    "application/json",
			},
		},
		{
			name:               fmt.Sprintf("%s JSON negative test #1", http.MethodPost),
			requestBody:        `{"link":"https://codepen.io"}`,
			requestPath:        "/api/shorten",
			requestMethod:      http.MethodPost,
			requestContentType: "",
			want: want{
				code: http.StatusUnsupportedMediaType,
				//contentType: "application/json",
			},
		},
		{
			name:               fmt.Sprintf("%s JSON negative test #2", http.MethodPost),
			requestBody:        `{"link":"https://codepen.io"}`,
			requestPath:        "/api/shorten",
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
		{
			name:                  fmt.Sprintf("%s GZIP positive test #1", http.MethodPost),
			requestBody:           `{"url":"https://codepen.io"}`,
			requestPath:           "/api/shorten",
			requestMethod:         http.MethodPost,
			requestContentType:    "application/json",
			requestAcceptEncoding: "gzip",
			want: want{
				code:            http.StatusCreated,
				contentType:     "application/json",
				contentEncoding: "gzip",
			},
		},
		{
			name:                   fmt.Sprintf("%s GZIP positive test #2", http.MethodPost),
			requestCompressBody:    utils.GzipCompress([]byte(`{"url":"https://codepen.io"}`)),
			requestPath:            "/api/shorten",
			requestMethod:          http.MethodPost,
			requestContentType:     "application/json",
			requestContentEncoding: "gzip",
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
	}

	Handler := *handlers.MyHandler(&cfg, database)
	ts := httptest.NewServer(&Handler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBoby := []byte(tt.requestBody)
			//var body = []byte(tt.requestBody)
			reqURL := tt.requestPath + tt.request
			var request *http.Request
			if len(tt.requestCompressBody) > 0 {
				request = httptest.NewRequest(tt.requestMethod, reqURL, bytes.NewBuffer(tt.requestCompressBody))
			} else {
				request = httptest.NewRequest(tt.requestMethod, reqURL, bytes.NewBuffer(reqBoby))
			}
			request.Header.Set("Content-Type", tt.requestContentType)
			if tt.requestAcceptEncoding != "" {
				request.Header.Set("Accept-Encoding", tt.requestAcceptEncoding)
			}
			if tt.requestContentEncoding != "" {
				request.Header.Set("Content-Encoding", tt.requestContentEncoding)
			}
			// создаём новый Recorder
			w := httptest.NewRecorder()
			Handler.ServeHTTP(w, request)
			resp := w.Result()
			// Проверяем StatusCode

			respStatusCode := resp.StatusCode
			wantStatusCode := tt.want.code
			assert.Equal(t, respStatusCode, wantStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, respStatusCode))

			// Проверяем Location
			resLocation := resp.Header.Get("Location")
			wantLocation := tt.want.location
			assert.Equal(t, resLocation, wantLocation, fmt.Errorf("expected Location %s, got %s", wantLocation, resLocation))

			// Проверяем Content-Type
			resContentType := resp.Header.Get("Content-Type")
			wantContentType := tt.want.contentType
			assert.Equal(t, resContentType, wantContentType, fmt.Errorf("expected Content-Type %s, got %s", wantContentType, resContentType))

			// получаем и проверяем тело запроса
			defer resp.Body.Close()
			resBodyTmp, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(string(resBodyTmp))
			resBody := utils.ValidateShortURL(string(resBodyTmp))
			wantBody := tt.want.responseFormat
			assert.Equal(t, resContentType, wantContentType, fmt.Errorf("expected BodyCheck %v, got %v", wantBody, resBody))

		})
	}
}
