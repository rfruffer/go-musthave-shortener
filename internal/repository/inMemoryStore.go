package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type URLEntry struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type InMemoryStore struct {
	mu    sync.RWMutex
	store map[string]URLEntry
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{store: make(map[string]URLEntry)}
}

func (s *InMemoryStore) Save(shortID, originalURL, uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[shortID] = URLEntry{UUID: uuid, ShortURL: shortID, OriginalURL: originalURL}
	return nil
}

func (s *InMemoryStore) Get(shortID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.store[shortID]
	if !ok {
		return "", errors.New("URL not found")
	}
	return entry.OriginalURL, nil
}

func (s *InMemoryStore) SaveToFile(path string) error {
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

func (s *InMemoryStore) LoadFromFile(path string) error {
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
		var entry URLEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}
		s.store[entry.ShortURL] = entry
	}
	return scanner.Err()
}
