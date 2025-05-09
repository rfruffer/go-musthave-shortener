package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
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
	repo := repository.NewInMemoryStore()
	service := services.NewURLService(repo)
	handler := handlers.NewURLHandler(service, "")

	r := chi.NewRouter()
	r.Get("/{id}", handler.GetShortURLHandler)
	r.Post("/", handler.CreateShortURLHandler)
	r.Post("/api/shorten", handler.CreateShortJSONURLHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid request", http.StatusUnauthorized)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid request", http.StatusUnauthorized)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	handler.SetResultHost(server.URL)

	tests := []struct {
		name      string
		method    string
		path      string
		body      string
		request   string
		want      want
		useSaveID bool
		saveID    bool
		isJSON    bool
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
			name:   "POST JSON CORRECT",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   `{"url": "https://practicum.yandex.ru"}`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				response:    "",
			},
			saveID: true,
			isJSON: true,
		},
		{
			name:   "GET CORRECT",
			method: http.MethodGet,
			path:   "/",
			body:   "",
			want: want{
				contentType: "text/html; charset=utf-8",
				statusCode:  http.StatusTemporaryRedirect,
				response:    "https://practicum.yandex.ru/",
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
				response:    "empty or invalid body\n",
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
				response:    "cant find id in store\n",
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
				response:    "invalid request\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			if tt.useSaveID {
				path = "/" + savedID
			}

			client := resty.NewWithClient(&http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse // <- ключевая строка
				},
			})

			req := client.R()

			if tt.isJSON {
				req.SetHeader("Content-Type", "application/json")
			} else {
				req.SetHeader("Content-Type", "text/plain")
			}

			if tt.method == http.MethodPost {
				req.SetBody(tt.body)
			}

			var resp *resty.Response
			var err error

			request := client.R().
				SetHeader("Content-Type", "text/plain")

			switch tt.method {
			case http.MethodPost:
				request.SetBody(tt.body)
				resp, err = request.Post(server.URL + path)
			case http.MethodGet:
				resp, err = request.Get(server.URL + path)
			default:
				resp, err = request.Execute(tt.method, server.URL+path)
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))

			body := string(resp.Body())

			if tt.saveID {
				if tt.isJSON == true {
					var jsonResp models.ShortenResponse
					err := json.Unmarshal(resp.Body(), &jsonResp)
					require.NoError(t, err)
					require.True(t, strings.HasPrefix(jsonResp.Result, server.URL+"/"))
				} else {
					savedID = strings.TrimPrefix(body, server.URL+"/")
					assert.Equal(t, server.URL+"/"+savedID, body)
				}
			} else if tt.useSaveID {
				assert.Equal(t, tt.want.response, resp.Header().Get("Location"))
			} else {
				assert.Equal(t, tt.want.response, body)
			}
		})
	}
}
