package repository

import "errors"

type InMemoryStore struct {
	store map[string]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{store: make(map[string]string)}
}

func (s *InMemoryStore) Save(shortID string, originalURL string) error {
	s.store[shortID] = originalURL
	return nil
}

func (s *InMemoryStore) Get(shortID string) (string, error) {
	originalURL, ok := s.store[shortID]
	if !ok {
		return "", errors.New("URl not found")
	}
	return originalURL, nil
}
