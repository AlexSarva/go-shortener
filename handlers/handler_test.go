package handlers_test

import (
	"AlexSarva/go-shortener/compressor"
	"AlexSarva/go-shortener/compressor/compress"
	"AlexSarva/go-shortener/constant"
	"AlexSarva/go-shortener/handlers"
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/utils"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
)

const ShortLen int = 10

type Compressor struct {
	compressor compressor.Compressor
}

func NewCompressor(data []byte) *Compressor {
	GzipCompressor := *compress.NewGzipCompress(data)
	return &Compressor{
		compressor: GzipCompressor,
	}
}

func TestMyHandler(t *testing.T) {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)
	database := app.NewStorage()
	CurCompressor := NewCompressor([]byte(`{"url":"https://codepen.io"}`))
	compressData := CurCompressor.compressor.Compress()
	insErr := database.Repo.InsertURL("Hasfe", "https://codepen.io", cfg.BaseURL, "ff2d2c4c-7bf7-49a7-a468-9c6d32aff40a")
	if insErr != nil {
		log.Println(insErr)
	}
	type want struct {
		location        string
		contentType     string
		contentEncoding string
		response        string
		code            int
		responseFormat  bool
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
			name:          fmt.Sprintf("%s ping test #1", http.MethodGet),
			requestMethod: http.MethodGet,
			requestPath:   "/ping",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          fmt.Sprintf("%s get many test #1", http.MethodGet),
			requestMethod: http.MethodGet,
			requestPath:   "/api/user/urls",
			want: want{
				code:        http.StatusNoContent,
				contentType: "application/json",
			},
		},
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
			name: fmt.Sprintf("%s JSON positive test #2", http.MethodPost),
			requestBody: `[
    {
        "correlation_id": "1",
        "original_url": "https://t.me/dvachannel/93381"
    },
    {
        "correlation_id": "2",
        "original_url": "https://t.me/moscowach/14075"
    }
]`,
			requestPath:        "/api/shorten/batch",
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
			requestCompressBody:    compressData,
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
	delCh := make(chan models.DeleteURL)
	Handler := *handlers.MyHandler(database, delCh)
	ts := httptest.NewServer(&Handler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBoby := []byte(tt.requestBody)
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
			resBody := utils.ValidateShortURL(string(resBodyTmp), cfg.BaseURL, ShortLen)
			wantBody := tt.want.responseFormat
			assert.Equal(t, resContentType, wantContentType, fmt.Errorf("expected BodyCheck %v, got %v", wantBody, resBody))

		})
	}
}
