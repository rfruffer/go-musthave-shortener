package router

import (
	"net/http"

	"github.com/gorilla/mux"
	// "github.com/rfruffer/go-musthave-shortener/config"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/middlewares"

	// "github.com/rfruffer/go-musthave-shortener/internal/repository"
	// "github.com/rfruffer/go-musthave-shortener/internal/services"
	"go.uber.org/zap"
)

type Router struct {
	UrlHandler *handlers.URLHandler
}

func SetupRouter(rt Router) http.Handler {
	r := mux.NewRouter()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	middlewares.InitLogger(sugar)
	r.Use(middlewares.LoggingMiddleware)

	r.HandleFunc("/{id}", rt.UrlHandler.GetShortURLHandler).Methods("GET")
	r.HandleFunc("/", rt.UrlHandler.CreateShortURLHandler).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.GzipMiddleware)
	api.HandleFunc("/shorten", rt.UrlHandler.CreateShortJSONURLHandler).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "invalid request", http.StatusBadRequest)
	})

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "invalid request", http.StatusBadRequest)
	})

	return r
}
