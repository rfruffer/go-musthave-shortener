package main

import (
	"net/http"

	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

func main() {
	urlService := services.NewUrlService()
	urlHandler := handlers.NewUrlHandler(urlService)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, urlHandler.ShortUrlHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
