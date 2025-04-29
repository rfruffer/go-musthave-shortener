package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rfruffer/go-musthave-shortener/config"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

func main() {
	cfg := config.ParseFlags()

	urlService := services.NewURLService()
	urlHandler := handlers.NewURLHandler(urlService, cfg.ResultHost)

	r := chi.NewRouter()

	r.Get("/{id}", urlHandler.GetShortURLHandler)
	r.Post("/", urlHandler.CreateShortURLHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid request", http.StatusBadRequest)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid request", http.StatusBadRequest)
	})

	err := http.ListenAndServe(cfg.StartHost, r)
	if err != nil {
		panic(err)
	}
}
