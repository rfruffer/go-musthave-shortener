package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
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

	var savedID string
	service := services.NewURLService()
	handler := handlers.NewURLHandler(service)

	r := chi.NewRouter()
	r.Get("/{id}", handler.GetShortURLHandler)
	r.Post("/", handler.CreateShortURLHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request", http.StatusUnauthorized)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request", http.StatusUnauthorized)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	tests := []struct {
		name      string
		method    string
		path      string
		body      string
		request   string
		want      want
		useSaveID bool
		saveID    bool
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
			saveID: true,
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
			useSaveID: true,
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
			path:   "/fakeId",
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "Cant find id in store\n",
			},
			useSaveID: false,
		},
		{
			name:   "INCORRECT METHOD",
			method: http.MethodPut,
			path:   "/",
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusUnauthorized,
				response:    "Invalid request\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			if tt.useSaveID {
				path = "/" + savedID
			}
			client := resty.New()

			resp, err := client.R().
				SetBody(tt.body).
				SetHeader("Content-Type", "text/plain").
				Execute(tt.method, server.URL+path)

			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))

			body := string(resp.Body())

			if tt.saveID {
				savedID = strings.TrimPrefix(body, "http://localhost:8080/")
				assert.Equal(t, tt.want.response+savedID, body)
			} else {
				assert.Equal(t, tt.want.response, body)
			}

		})
	}
}
