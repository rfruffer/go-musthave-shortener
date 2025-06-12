package tests

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/middlewares"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
	"github.com/stretchr/testify/require"
)

func TestGzipCompression(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInFileStore()
	service := services.NewURLService(repo)
	handler := handlers.NewURLHandler(service, "")

	r := gin.New()
	r.Use(middlewares.GinGzipMiddleware())
	r.POST("/api/shorten", handler.CreateShortJSONURLHandler)

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusUnauthorized, "invalid request")
	})
	r.NoMethod(func(c *gin.Context) {
		c.String(http.StatusUnauthorized, "invalid request")
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
		require.NoError(t, zb.Close())

		req, err := http.NewRequest("POST", server.URL+"/api/shorten", buf)
		require.NoError(t, err)

		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		defer zr.Close()

		unzipped, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.Contains(t, string(unzipped), `"result":"`+server.URL+`/`)
	})
}
