package main

import (
	"net/http"

	// "github.com/go-chi/chi"
	"github.com/gorilla/mux"
	"github.com/rfruffer/go-musthave-shortener/config"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/middlewares"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
	"go.uber.org/zap"
)

func main() {
	cfg := config.ParseFlags()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	repo := repository.NewInMemoryStore()
	service := services.NewURLService(repo)
	handler := handlers.NewURLHandler(service, cfg.ResultHost)

	// r := chi.NewRouter()
	r := mux.NewRouter()

	middlewares.InitLogger(sugar)
	r.Use(middlewares.LoggingMiddleware)
	api := r.PathPrefix("/").Subrouter()
	api.Use(middlewares.GzipMiddleware)

	r.HandleFunc("/{id}", handler.GetShortURLHandler).Methods("Get")
	// r.Get("/{id}", handler.GetShortURLHandler)
	// r.Post("/", handler.CreateShortURLHandler)
	r.HandleFunc("/", handler.CreateShortURLHandler).Methods("Post")
	// api.Post("/api/shorten", handler.CreateShortJSONURLHandler)
	api.HandleFunc("/api/shorten", handler.CreateShortJSONURLHandler).Methods("Post")

	// r.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// 	http.Error(w, "invalid request", http.StatusBadRequest)
	// })

	// r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
	// 	http.Error(w, "invalid request", http.StatusBadRequest)
	// })

	err = http.ListenAndServe(cfg.StartHost, r)
	if err != nil {
		panic(err)
	}
}
