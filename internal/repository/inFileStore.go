package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/rfruffer/go-musthave-shortener/internal/models"
)

// URLHandler предоставляет методы для работы с файлами.
type InFileStore struct {
	mu    sync.RWMutex
	store map[string]models.URLEntry
}

// NewInFileStore создает новый InFileStore
func NewInFileStore() *InFileStore {
	return &InFileStore{store: make(map[string]models.URLEntry)}
}

// Save сохраняет файл
func (s *InFileStore) Save(shortID, originalURL, uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[shortID] = models.URLEntry{UUID: uuid, ShortURL: shortID, OriginalURL: originalURL}
	return nil
}

// GetURLByShort получить URL по короткой ссылке
func (s *InFileStore) GetURLByShort(shortID string) (models.URLEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.store[shortID]
	if !ok {
		return models.URLEntry{}, errors.New("URL not found")
	}
	return entry, nil
}

// SaveToFile сохраняет файл
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

// LoadFromFile загрузить из файла
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

// Ping заглушка из интерфейса
func (s *InFileStore) Ping() error {
	return nil
}

// Ping заглушка из интерфейса
func (s *InFileStore) GetShortIDByOriginalURL(originalURL string) (string, error) {
	return "", nil
}

// Ping заглушка из интерфейса
func (s *InFileStore) GetByUser(userID string) ([]models.URLEntry, error) {
	return []models.URLEntry{}, nil
}

// Ping заглушка из интерфейса
func (s *InFileStore) MarkURLsDeleted(userID string, ids []string) error {
	return nil
}
