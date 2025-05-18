package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer io.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func GinGzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Decompress incoming request
		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") &&
			strings.HasPrefix(c.Request.Header.Get("Content-Type"), "application/json") &&
			c.Request.Method == http.MethodPost {

			gr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to decompress request"})
				return
			}
			defer gr.Close()

			c.Request.Body = io.NopCloser(gr)
		}

		// Compress outgoing response
		if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")

			gzw := gzip.NewWriter(c.Writer)
			defer gzw.Close()

			c.Writer = &gzipWriter{
				ResponseWriter: c.Writer,
				writer:         gzw,
			}

			c.Next()
			_ = gzw.Close()
		} else {
			c.Next()
		}
	}
}
