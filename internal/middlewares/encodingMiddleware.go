package middlewares

import (
	"compress/gzip"
	"fmt"
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
		fmt.Println("[DEBUG] FullPath:", c.FullPath(), "URI:", c.Request.RequestURI)

		if !strings.HasPrefix(c.Request.URL.Path, "/api/") && c.Request.URL.Path != "/api" {
			c.Next()
			return
		}
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			defer gr.Close()
			c.Request.Body = io.NopCloser(gr)
		}

		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")
		if contentType == "" {
			contentType = http.DetectContentType([]byte{})
		}

		if strings.Contains(contentType, "application/json") ||
			strings.Contains(contentType, "text/html") ||
			strings.Contains(contentType, "text/plain") {

			gzw := gzip.NewWriter(c.Writer)
			defer gzw.Close()

			c.Header("Content-Encoding", "gzip")
			// c.Header("Vary", "Accept-Encoding")

			c.Writer = &gzipWriter{
				ResponseWriter: c.Writer,
				writer:         gzw,
			}

			c.Next()
			_ = gzw.Close()
			return
		}
		c.Next()
	}
}
