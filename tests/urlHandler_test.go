package tests

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUrlHandler_ShortUrlHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}

	var savedId string
	service := services.NewUrlService()
	handler := handlers.NewUrlHandler(service)

	tests := []struct {
		name      string
		method    string
		path      string
		body      string
		request   string
		want      want
		useSaveId bool
		saveId    bool
	}{
		{
			name:   "POST CORRECT",
			method: http.MethodPost,
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				response:    "http://localhost:8080/",
			},
			saveId: true,
		},
		{
			name:   "GET CORRECT",
			method: http.MethodGet,
			path:   "/",
			body:   "",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusTemporaryRedirect,
				response:    "Location: https://practicum.yandex.ru/",
			},
			useSaveId: true,
		},
		{
			name:   "NEGATIVE POST",
			method: http.MethodPost,
			path:   "/",
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "Empty or invalid body\n",
			},
		},
		{
			name:   "NEGATIVE GET",
			method: http.MethodGet,
			path:   "/",
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "Missing ID\n",
			},
			useSaveId: false,
		},
		{
			name:   "INCORRECT METHOD",
			method: http.MethodPut,
			path:   "/",
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusUnauthorized,
				response:    "Unsupport method\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			if tt.useSaveId {
				path = "/" + savedId
			}

			r := httptest.NewRequest(tt.method, path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			http.HandlerFunc(handler.ShortUrlHandler).ServeHTTP(w, r)

			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if tt.saveId {
				savedId = strings.TrimPrefix(string(bodyResult), "http://localhost:8080/")
				assert.Equal(t, tt.want.response+savedId, string(bodyResult))
			} else {
				assert.Equal(t, tt.want.response, string(bodyResult))
			}

		})
	}
}
