package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
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
	if c.w.Header().Get("Content-Encoding") != "gzip" {
		return c.w.Write(p)
	}
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode >= 300 && statusCode < 400 {
		c.w.WriteHeader(statusCode)
		return
	}
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
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

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") &&
			strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") &&
			r.Method == http.MethodPost {

			cr, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "failed to decompress request", http.StatusBadRequest)
				return
			}
			defer cr.Close()
			r.Body = cr
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		cw := newCompressWriter(w)
		defer cw.Close()

		h.ServeHTTP(cw, r)
	})
}
