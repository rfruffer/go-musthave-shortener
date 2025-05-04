package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func InitLogger(logger *zap.SugaredLogger) {
	sugar = *logger
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respData := &responseData{}
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   respData,
		}

		h.ServeHTTP(lw, r)

		duration := time.Since(start)

		sugar.Infow("incoming request",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", respData.status,
			"duration", duration,
			"size", respData.size,
		)
	})
}
