package tests

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/middlewares"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
	"github.com/stretchr/testify/require"
)

func TestGzipCompression(t *testing.T) {
	repo := repository.NewInMemoryStore()
	service := services.NewURLService(repo)
	handler := handlers.NewURLHandler(service, "")

	r := chi.NewRouter()
	r.Use(middlewares.GzipMiddleware)

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

	requestBody := `{"url": "https://practicum.yandex.ru/"}`

	t.Run("sends_gzip", func(t *testing.T) {

		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", server.URL+"/api/shorten", buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		defer zr.Close()

		unzipped, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.Contains(t, string(unzipped), `"result":"`+server.URL+`/`)
	})
}
