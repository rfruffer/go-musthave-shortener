package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	// "sync"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
)

var shortSize = 8

// URLService подключает интерфейс StoreRepositoryInterface
type URLService struct {
	repo repository.StoreRepositoryInterface
}

// NewURLService создает новый URLService
func NewURLService(repo repository.StoreRepositoryInterface) *URLService {
	return &URLService{repo: repo}
}

// GenerateShortURL функция генерации короткиз ссылок
func (s *URLService) GenerateShortURL(originalURL string, uuid string) (string, error) {
	b := make([]byte, shortSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	id := base64.URLEncoding.EncodeToString(b)[:shortSize]
	// uuid := uuid.New().String()

	err = s.repo.Save(id, originalURL, uuid)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			existingID, getErr := s.repo.GetShortIDByOriginalURL(originalURL)
			if getErr != nil {
				return "", getErr
			}
			return existingID, repository.ErrAlreadyExists
		}
		return "", err
	}

	return id, nil
}

// RedirectURL производит редирект на указанный адрес
func (s *URLService) RedirectURL(id string) (string, error) {
	URL, err := s.repo.GetURLByShort(id)
	if err != nil {
		return "", fmt.Errorf("cant find id in store")
	}
	if URL.DeletedFlag {
		return "", repository.ErrGone
	}
	return URL.OriginalURL, nil
}

// Ping проверка подключения к базе
func (s *URLService) Ping() error {
	return s.repo.Ping()
}

// GenerateBatchShortURLs генерация коротких урл при пакетной загрузке
func (s *URLService) GenerateBatchShortURLs(req []models.BatchOriginalURL, userID string) ([]models.BatchShortURL, error) {
	resp := make([]models.BatchShortURL, 0, len(req))

	for _, item := range req {
		id, err := s.GenerateShortURL(item.OriginalURL, userID)
		if err != nil {
			return nil, err
		}
		resp = append(resp, models.BatchShortURL{
			CorrelationID: item.CorrelationID,
			ShortURL:      id,
		})
	}
	return resp, nil
}

// GetURLsByUser получить все ссылки загруженные пользователем
func (s *URLService) GetURLsByUser(userID string) ([]models.URLEntry, error) {
	return s.repo.GetByUser(userID)
}

// DeleteUserURLs удалить ссылки пользователя
func (s *URLService) DeleteUserURLs(userID string, ids []string) error {
	return s.repo.MarkURLsDeleted(userID, ids)
}
