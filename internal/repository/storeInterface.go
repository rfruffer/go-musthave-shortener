package repository

import "github.com/rfruffer/go-musthave-shortener/internal/models"

type StoreRepositoryInterface interface {
	Save(shortID string, originalURL string, uuid string) error
	Get(shortID string) (string, error)

	SaveToFile(path string) error
	LoadFromFile(path string) error

	Ping() error
	GetShortIDByOriginalURL(originalURL string) (string, error)
	GetByUser(userID string) ([]models.URLEntry, error)
}
