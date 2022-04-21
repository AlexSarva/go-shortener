package handlers_test

import (
	"AlexSarva/go-shortener/checker"
	"AlexSarva/go-shortener/handlers"
	"AlexSarva/go-shortener/internal/app"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeHandlerGet(t *testing.T) {
	var database = app.InitDB()
	database.Insert("https://codepen.io", "Hasfe")
	type want struct {
		code     int
		location string
		//response    string
		contentType string
	}

	tests := []struct {
		request string
		name    string
		want    want
	}{
		// TODO: Add test cases.
		{
			name:    "positive test #1",
			request: "Hasfe",
			want: want{
				code:        http.StatusTemporaryRedirect,
				location:    "https://codepen.io",
				contentType: "text/plain",
			},
		},
		{
			name:    "negative test #1",
			request: "",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:    "negative test #2",
			request: "sometext",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.request, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := handlers.MakeHandler(database)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()
			// проверяем код ответа
			resStatusCode := res.StatusCode
			wantStatusCode := tt.want.code
			assert.Equal(t, resStatusCode, wantStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, resStatusCode))

			// Проверяем Location
			resLocation := res.Header.Get("Location")
			wantLocation := tt.want.location
			assert.Equal(t, resLocation, wantLocation, fmt.Errorf("expected StatusCode %s, got %s", wantLocation, resLocation))

			// Проверяем Content-Type
			resContentType := res.Header.Get("Content-Type")
			wantContentType := tt.want.contentType
			assert.Equal(t, resContentType, wantContentType, fmt.Errorf("expected Content-Type %s, got %s", wantContentType, resContentType))

			// получаем и проверяем тело запроса
			//defer res.Body.Close()
			//resBody, err := io.ReadAll(res.Body)
			//if err != nil {
			//	t.Fatal(err)
			//}
			//if string(resBody) != tt.want.response {
			//	t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			//}

		})
	}
}

func TestMakeHandlerPost(t *testing.T) {
	var database = app.InitDB()
	//database.Insert("https://codepen.io", "Hasfe")
	type want struct {
		code        int
		response    bool
		contentType string
	}

	tests := []struct {
		request string
		name    string
		want    want
	}{
		// TODO: Add test cases.
		{
			name:    "positive test #1",
			request: "https://codepen.io",
			want: want{
				code:        http.StatusCreated,
				response:    true,
				contentType: "text/plain",
			},
		},
		{
			name:    "positive test #2",
			request: "www.google.com",
			want: want{
				code:        http.StatusCreated,
				response:    true,
				contentType: "text/plain",
			},
		},
		{
			name:    "positive test #3",
			request: "www.url-with-querystring.com/?url=has-querystring",
			want: want{
				code:        http.StatusCreated,
				response:    true,
				contentType: "text/plain",
			},
		},
		{
			name:    "negative test #1",
			request: "https://codepen",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:    "negative test #2",
			request: "что-то написано",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:    "negative test #3",
			request: "",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body = []byte(tt.request)
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := handlers.MakeHandler(database)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()
			// проверяем код ответа
			resStatusCode := res.StatusCode
			wantStatusCode := tt.want.code
			assert.Equal(t, resStatusCode, wantStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, resStatusCode))

			// Проверяем Content-Type
			resContentType := res.Header.Get("Content-Type")
			wantContentType := tt.want.contentType
			assert.Equal(t, resContentType, wantContentType, fmt.Errorf("expected Content-Type %s, got %s", wantContentType, resContentType))

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBodyTmp, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			resBody := checker.CheckShortUrl(string(resBodyTmp))
			wantBody := tt.want.response
			assert.Equal(t, resContentType, wantContentType, fmt.Errorf("expected BodyCheck %v, got %v", wantBody, resBody))
		})
	}
}
