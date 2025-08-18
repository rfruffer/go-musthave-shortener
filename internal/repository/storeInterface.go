package repository

import "github.com/rfruffer/go-musthave-shortener/internal/models"

// StoreRepositoryInterface представляет методы для взаимодействия с хранилищем
type StoreRepositoryInterface interface {
	Save(shortID string, originalURL string, uuid string) error
	GetURLByShort(shortID string) (models.URLEntry, error)

	SaveToFile(path string) error
	LoadFromFile(path string) error

	Ping() error
	GetShortIDByOriginalURL(originalURL string) (string, error)
	GetByUser(userID string) ([]models.URLEntry, error)
	MarkURLsDeleted(userID string, ids []string) error
}
