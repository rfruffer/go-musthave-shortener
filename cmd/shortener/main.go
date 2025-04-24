package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

func main() {
	urlService := services.NewURLService()
	urlHandler := handlers.NewURLHandler(urlService)

	r := chi.NewRouter()

	r.Get("/{id}", urlHandler.GetShortURLHandler)
	r.Post("/", urlHandler.CreateShortURLHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	})

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
