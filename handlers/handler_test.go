package handlers_test

import (
	"AlexSarva/go-shortener/handlers"
	"AlexSarva/go-shortener/storage"
	"AlexSarva/go-shortener/utils"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMytHandler(t *testing.T) {
	database := *storage.InitDB()
	insErr := database.Insert("https://codepen.io", "Hasfe")
	if insErr != nil {
		log.Println(insErr)
	}
	type want struct {
		code           int
		location       string
		contentType    string
		responseFormat bool
		response       string
	}

	tests := []struct {
		name          string
		request       string
		requestMethod string
		requestBody   string
		want          want
	}{
		{
			name:          fmt.Sprintf("%s positive test #1", http.MethodGet),
			request:       "Hasfe",
			requestMethod: http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				location:    "https://codepen.io",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodGet),
			request:       "",
			requestMethod: http.MethodGet,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #2", http.MethodGet),
			request:       "sometext",
			requestMethod: http.MethodGet,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s positive test #1", http.MethodPost),
			requestBody:   "https://codepen.io",
			requestMethod: http.MethodPost,
			want: want{
				code:           http.StatusCreated,
				responseFormat: true,
				contentType:    "text/plain; charset=utf-8",
			},
		},
		{
			name:          fmt.Sprintf("%s positive test #2", http.MethodPost),
			requestBody:   "www.google.com",
			requestMethod: http.MethodPost,
			want: want{
				code:           http.StatusCreated,
				responseFormat: true,
				contentType:    "text/plain; charset=utf-8",
			},
		},
		{
			name:          fmt.Sprintf("%s positive test #3", http.MethodPost),
			requestBody:   "www.url-with-querystring.com/?url=has-querystring",
			requestMethod: http.MethodPost,
			want: want{
				code:           http.StatusCreated,
				responseFormat: true,
				contentType:    "text/plain; charset=utf-8",
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodPost),
			requestBody:   "https://codepen",
			requestMethod: http.MethodPost,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #2", http.MethodPost),
			requestBody:   "something",
			requestMethod: http.MethodPost,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #3", http.MethodPost),
			requestBody:   "",
			requestMethod: http.MethodPost,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodPut),
			requestBody:   "something",
			request:       "something",
			requestMethod: http.MethodPut,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodConnect),
			requestBody:   "something",
			request:       "something",
			requestMethod: http.MethodConnect,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodDelete),
			requestBody:   "something",
			request:       "something",
			requestMethod: http.MethodDelete,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodHead),
			requestBody:   "something",
			request:       "something",
			requestMethod: http.MethodHead,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodOptions),
			requestBody:   "something",
			request:       "something",
			requestMethod: http.MethodOptions,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodPatch),
			requestBody:   "something",
			request:       "something",
			requestMethod: http.MethodPatch,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", http.MethodTrace),
			requestBody:   "something",
			request:       "something",
			requestMethod: http.MethodTrace,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s negative test #1", "CUSTOM"),
			requestBody:   "something",
			request:       "something",
			requestMethod: "CUSTOM",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	Handler := *handlers.MyHandler(database)
	ts := httptest.NewServer(&Handler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBoby := []byte(tt.requestBody)
			//var body = []byte(tt.requestBody)
			reqUrl := "/" + tt.request
			request := httptest.NewRequest(tt.requestMethod, reqUrl, bytes.NewBuffer(reqBoby))
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
			resBody := utils.ValidateShortUrl(string(resBodyTmp))
			wantBody := tt.want.responseFormat
			assert.Equal(t, resContentType, wantContentType, fmt.Errorf("expected BodyCheck %v, got %v", wantBody, resBody))

		})
	}
}
