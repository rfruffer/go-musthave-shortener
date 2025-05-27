package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rfruffer/go-musthave-shortener/cmd/shortener/router"
	"github.com/rfruffer/go-musthave-shortener/config"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	posgreConfig "github.com/rfruffer/go-musthave-shortener/internal/repository/posgreConfig"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

func main() {
	cfg := config.ParseFlags()

	var repo repository.StoreRepositoryInterface
	var service *services.URLService
	var shortURLHandler *handlers.URLHandler

	switch cfg.Storage {
	case "postgres":
		db, err := posgreConfig.InitDB(cfg.DBDSN)
		if err != nil {
			log.Fatalf("failed to initialize database: %v", err)
		}
		defer posgreConfig.CloseDB(db)
		repo = repository.NewDBStore(db)

		service = services.NewURLService(repo)
		shortURLHandler = handlers.NewURLHandler(service, cfg.ResultHost, repo)
	default:
		repo = repository.NewInFileStore()
		service = services.NewURLService(repo)
		shortURLHandler = handlers.NewURLHandler(service, cfg.ResultHost, nil)
	}

	if err := repo.LoadFromFile(cfg.FilePath); err != nil {
		log.Fatalf("failed to load from file: %v", err)
	}

	r := router.SetupRouter(router.Router{
		URLHandler: shortURLHandler,
	})

	server := &http.Server{
		Addr:    cfg.StartHost,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("starting server on %s", cfg.StartHost)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	<-stop
	log.Println("shutting down server...")

	if err := server.Close(); err != nil {
		log.Printf("error shutting down server: %v", err)
	}

	if err := repo.SaveToFile(cfg.FilePath); err != nil {
		log.Printf("failed to save to file: %v", err)
	}

	log.Println("server stopped gracefully")
}
