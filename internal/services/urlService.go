package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
)

var shortSize = 8

type URLService struct {
	repo repository.StoreRepositoryInterface
}

func NewURLService(repo repository.StoreRepositoryInterface) *URLService {
	return &URLService{repo: repo}
}

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

func (s *URLService) RedirectURL(id string) (string, error) {
	originalURL, err := s.repo.Get(id)
	if err != nil {
		return "", fmt.Errorf("cant find id in store")
	}
	return originalURL, nil
}

func (s *URLService) Ping() error {
	return s.repo.Ping()
}

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

func (s *URLService) GetURLsByUser(userID string) ([]models.URLEntry, error) {
	return s.repo.GetByUser(userID)
}
