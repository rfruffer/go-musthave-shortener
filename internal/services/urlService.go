package services

import (
	"encoding/base64"
	"fmt"
	"math/rand"
)

var (
	store     = make(map[string]string)
	shortSize = 8
)

type UrlService struct {
}

func NewUrlService() *UrlService {
	return &UrlService{}
}

func (s *UrlService) GenerateShortUrl(originalURL string) (string, error) {
	b := make([]byte, shortSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	id := base64.URLEncoding.EncodeToString(b)[:shortSize]
	store[id] = originalURL
	return id, nil
}

func (s *UrlService) RedirectUrl(id string) (string, error) {
	originalURL, ok := store[id]

	if !ok {
		return "", fmt.Errorf("Cant find id in store")
	}
	return originalURL, nil
}
