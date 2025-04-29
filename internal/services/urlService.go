package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

var (
	store     = make(map[string]string)
	shortSize = 8
)

type URLService struct {
}

func NewURLService() *URLService {
	return &URLService{}
}

func (s *URLService) GenerateShortURL(originalURL string) (string, error) {
	b := make([]byte, shortSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	id := base64.URLEncoding.EncodeToString(b)[:shortSize]
	store[id] = originalURL
	return id, nil
}

func (s *URLService) RedirectURL(id string) (string, error) {
	originalURL, ok := store[id]

	if !ok {
		return "", fmt.Errorf("cant find id in store")
	}
	return originalURL, nil
}
