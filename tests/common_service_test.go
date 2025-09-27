package tests

import (
	"testing"

	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

func TestCommonServiceCreateShortURL(t *testing.T) {
	// Создаем test environment
	repo := repository.NewInFileStore()
	urlService := services.NewURLService(repo)
	commonService := services.NewCommonURLService(urlService, "http://localhost:8080")

	// Тест создания короткой ссылки
	originalURL := "https://example.com"
	userID := "test-user-123"

	shortURL, alreadyExists, err := commonService.CreateShortURL(originalURL, userID)
	if err != nil {
		t.Fatalf("CreateShortURL failed: %v", err)
	}

	if shortURL == "" {
		t.Error("Expected non-empty short URL")
	}

	if alreadyExists {
		t.Error("Expected AlreadyExists to be false for new URL")
	}

	t.Logf("Created short URL: %s", shortURL)

	// Тест повторного создания той же ссылки
	shortURL2, alreadyExists2, err2 := commonService.CreateShortURL(originalURL, userID)
	if err2 != nil {
		t.Fatalf("Second CreateShortURL failed: %v", err2)
	}

	if !alreadyExists2 {
		t.Error("Expected AlreadyExists to be true for existing URL")
	}

	if shortURL != shortURL2 {
		t.Errorf("Expected same short URL, got %s and %s", shortURL, shortURL2)
	}
}

func TestCommonServiceBatch(t *testing.T) {
	// Создаем test environment
	repo := repository.NewInFileStore()
	urlService := services.NewURLService(repo)
	commonService := services.NewCommonURLService(urlService, "http://localhost:8080")

	// Тест пакетного создания
	requests := []models.BatchOriginalURL{
		{
			CorrelationID: "1",
			OriginalURL:   "https://example1.com",
		},
		{
			CorrelationID: "2",
			OriginalURL:   "https://example2.com",
		},
	}

	userID := "test-user-batch"

	results, err := commonService.BatchCreateShortURLs(requests, userID)
	if err != nil {
		t.Fatalf("BatchCreateShortURLs failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	for i, result := range results {
		if result.ShortURL == "" {
			t.Errorf("Expected non-empty short URL for item %d", i)
		}
		if result.CorrelationID == "" {
			t.Errorf("Expected non-empty correlation ID for item %d", i)
		}
		t.Logf("Batch item %d: %s -> %s", i, result.CorrelationID, result.ShortURL)
	}
}

func TestCommonServicePing(t *testing.T) {
	// Создаем test environment
	repo := repository.NewInFileStore()
	urlService := services.NewURLService(repo)
	commonService := services.NewCommonURLService(urlService, "http://localhost:8080")

	// Тест ping
	err := commonService.Ping()
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	t.Log("Ping successful")
}

func TestCommonServiceGetStats(t *testing.T) {
	// Создаем test environment
	repo := repository.NewInFileStore()
	urlService := services.NewURLService(repo)
	commonService := services.NewCommonURLService(urlService, "http://localhost:8080")

	// Создаем несколько URL для статистики
	commonService.CreateShortURL("https://example1.com", "user1")
	commonService.CreateShortURL("https://example2.com", "user2")

	// Тест получения статистики
	urlCount, userCount, err := commonService.GetStats()
	if err != nil {
		t.Errorf("GetStats failed: %v", err)
	}

	if urlCount < 2 {
		t.Errorf("Expected at least 2 URLs, got %d", urlCount)
	}

	if userCount < 2 {
		t.Errorf("Expected at least 2 users, got %d", userCount)
	}

	t.Logf("Stats: %d URLs, %d users", urlCount, userCount)
}