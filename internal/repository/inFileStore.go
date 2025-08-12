package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/rfruffer/go-musthave-shortener/internal/models"
)

type InFileStore struct {
	mu    sync.RWMutex
	store map[string]models.URLEntry
}

func NewInFileStore() *InFileStore {
	return &InFileStore{store: make(map[string]models.URLEntry)}
}

func (s *InFileStore) Save(shortID, originalURL, uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[shortID] = models.URLEntry{UUID: uuid, ShortURL: shortID, OriginalURL: originalURL}
	return nil
}

func (s *InFileStore) GetURLByShort(shortID string) (models.URLEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.store[shortID]
	if !ok {
		return models.URLEntry{}, errors.New("URL not found")
	}
	return entry, nil
}

func (s *InFileStore) SaveToFile(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	for _, entry := range s.store {
		if err := enc.Encode(entry); err != nil {
			return err
		}
	}
	return nil
}

func (s *InFileStore) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry models.URLEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}
		s.store[entry.ShortURL] = entry
	}
	return scanner.Err()
}

func (s *InFileStore) Ping() error {
	return nil
}

func (s *InFileStore) GetShortIDByOriginalURL(originalURL string) (string, error) {
	return "", nil
}

func (s *InFileStore) GetByUser(userID string) ([]models.URLEntry, error) {
	return []models.URLEntry{}, nil
}

func (s *InFileStore) MarkURLsDeleted(userID string, ids []string) error {
	return nil
}
