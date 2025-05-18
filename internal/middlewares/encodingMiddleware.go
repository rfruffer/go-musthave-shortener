package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

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

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") &&
			strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") &&
			r.Method == http.MethodPost {

			cr, err := newCompressReader(r.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to decompress request"})
				return
			}
			defer cr.Close()
			r.Body = cr
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && r.Method == http.MethodPost {
			cw := newCompressWriter(w)
			defer cw.Close()

			c.Writer = &ginResponseWriter{
				ResponseWriter: w,
				writer:         cw.zw,
			}
			c.Header("Content-Encoding", "gzip")
		}

		c.Next()
	}
}
