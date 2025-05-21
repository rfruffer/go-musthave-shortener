package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
)

var (
	shortSize = 8
)

type URLService struct {
	repo *repository.InMemoryStore
}

func NewURLService(repo *repository.InMemoryStore) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) GenerateShortURL(originalURL string) (string, error) {
	b := make([]byte, shortSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	id := base64.URLEncoding.EncodeToString(b)[:shortSize]
	uuid := uuid.New().String()
	err = s.repo.Save(id, originalURL, uuid)
	if err != nil {
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
