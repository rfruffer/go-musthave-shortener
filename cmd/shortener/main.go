package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rfruffer/go-musthave-shortener/cmd/shortener/router"
	"github.com/rfruffer/go-musthave-shortener/config"
	"github.com/rfruffer/go-musthave-shortener/internal/async"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	posgreConfig "github.com/rfruffer/go-musthave-shortener/internal/repository/posgreConfig"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

// билд-переменные (заполняются через -ldflags)
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	// печать информации о сборке (stdout)
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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
		shortURLHandler = handlers.NewURLHandler(service, cfg.ResultHost)

		doneCh := make(chan struct{})
		queue1 := make(chan async.DeleteTask)

		merged := async.FanIn(doneCh, queue1)
		async.StartDeleteWorker(doneCh, repo, merged)
		shortURLHandler.DeleteChan = queue1

	default:
		repo = repository.NewInFileStore()
		service = services.NewURLService(repo)
		shortURLHandler = handlers.NewURLHandler(service, cfg.ResultHost)
	}

	if err := repo.LoadFromFile(cfg.FilePath); err != nil {
		log.Fatalf("failed to load from file: %v", err)
	}

	r := router.SetupRouter(router.Router{
		URLHandler: shortURLHandler,
		SecretKey:  cfg.SecretKey,
	})

	server := &http.Server{
		Addr:    cfg.StartHost,
		Handler: r,
	}

	// Создаем контекст, который отменяется при получении сигнала завершения
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	serverErr := make(chan error, 1)
	
	go func() {
		if cfg.EnableHTTPS {
			log.Printf("starting HTTPS server on %s", cfg.StartHost)
			if err := server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile); err != nil && err != http.ErrServerClosed {
				serverErr <- fmt.Errorf("HTTPS server error: %w", err)
				return
			}
		} else {
			log.Printf("starting HTTP server on %s", cfg.StartHost)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				serverErr <- fmt.Errorf("HTTP server error: %w", err)
				return
			}
		}
	}()

	// Ждем либо сигнал завершения, либо ошибку сервера
	select {
	case err := <-serverErr:
		log.Printf("server failed to start: %v", err)
	case <-ctx.Done():
		log.Println("shutting down server...")
	}

	// Создаем context с timeout для graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Закрываем асинхронные воркеры если они есть
	switch cfg.Storage {
	case "postgres":
		if shortURLHandler.DeleteChan != nil {
			log.Println("closing delete workers...")
			close(shortURLHandler.DeleteChan)
		}
	}

	// Graceful shutdown сервера - ждем завершения всех активных запросов
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("error during server shutdown: %v", err)
	} else {
		log.Println("all active requests completed")
	}

	// Сохраняем все несохраненные данные
	log.Println("saving data to storage...")
	if err := repo.SaveToFile(cfg.FilePath); err != nil {
		log.Printf("failed to save to file: %v", err)
	} else {
		log.Println("data saved successfully")
	}

	log.Println("server stopped gracefully")
}
