package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ginResponseWriter struct {
	gin.ResponseWriter
	writer io.Writer
}

func (g *ginResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func GinGzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		w := c.Writer

		// Декодируем входящий gzip-запрос
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") &&
			strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") &&
			r.Method == http.MethodPost {

			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to decompress request"})
				return
			}
			defer gr.Close()
			r.Body = io.NopCloser(gr)
		}

		// Сжимаем исходящий ответ, если клиент поддерживает gzip
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") &&
			strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {

			gzw := gzip.NewWriter(w)
			defer gzw.Close()

			c.Header("Content-Encoding", "gzip")
			c.Writer = &ginResponseWriter{
				ResponseWriter: w,
				writer:         gzw,
			}
		}

		c.Next()
	}
}
