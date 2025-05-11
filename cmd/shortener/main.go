package main

import (
	"log"
	"net/http"

	// "github.com/go-chi/chi"
	// "github.com/gorilla/mux"
	"github.com/rfruffer/go-musthave-shortener/cmd/shortener/router"
	"github.com/rfruffer/go-musthave-shortener/config"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"

	// "github.com/rfruffer/go-musthave-shortener/internal/middlewares"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
	// "go.uber.org/zap"
)

func main() {

	cfg := config.ParseFlags()
	repo := repository.NewInMemoryStore()
	service := services.NewURLService(repo)
	shortURLhandler := handlers.NewURLHandler(service, cfg.ResultHost)

	router := router.SetupRouter(router.Router{
		UrlHandler: shortURLhandler,
	})

	server := &http.Server{
		Addr:    cfg.StartHost,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}
