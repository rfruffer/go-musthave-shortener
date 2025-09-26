// Package services содержит общую бизнес-логику для HTTP и gRPC интерфейсов.
package services

import (
	"errors"

	"github.com/rfruffer/go-musthave-shortener/internal/async"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
)

// CommonURLService объединяет всю бизнес-логику для работы с сокращенными URL
// Используется как HTTP, так и gRPC хендлерами
type CommonURLService struct {
	urlService *URLService
	baseURL    string
	deleteChan chan async.DeleteTask
}

// NewCommonURLService создает новый CommonURLService
func NewCommonURLService(urlService *URLService, baseURL string) *CommonURLService {
	return &CommonURLService{
		urlService: urlService,
		baseURL:    baseURL,
	}
}

// SetDeleteChannel устанавливает канал для асинхронного удаления
func (s *CommonURLService) SetDeleteChannel(deleteChan chan async.DeleteTask) {
	s.deleteChan = deleteChan
}

// CreateShortURL создает короткую ссылку и возвращает результат с информацией о том, существовала ли ссылка ранее
func (s *CommonURLService) CreateShortURL(originalURL, userID string) (shortURL string, alreadyExists bool, err error) {
	if originalURL == "" {
		return "", false, errors.New("URL не может быть пустым")
	}

	id, err := s.urlService.GenerateShortURL(originalURL, userID)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return s.baseURL + "/" + id, true, nil
		}
		return "", false, err
	}

	return s.baseURL + "/" + id, false, nil
}

// GetOriginalURL получает оригинальный URL по короткому ID
func (s *CommonURLService) GetOriginalURL(id string) (originalURL string, deleted bool, err error) {
	if id == "" {
		return "", false, errors.New("ID не может быть пустым")
	}

	urlEntry, err := s.urlService.repo.GetURLByShort(id)
	if err != nil {
		return "", false, errors.New("URL не найден")
	}

	if urlEntry.DeletedFlag {
		return "", true, nil
	}

	return urlEntry.OriginalURL, false, nil
}

// Ping проверяет подключение к базе данных
func (s *CommonURLService) Ping() error {
	return s.urlService.Ping()
}

// BatchCreateShortURLs создает несколько коротких ссылок за один запрос
func (s *CommonURLService) BatchCreateShortURLs(requests []models.BatchOriginalURL, userID string) ([]models.BatchShortURL, error) {
	if len(requests) == 0 {
		return nil, errors.New("пустой список запросов")
	}

	resp := make([]models.BatchShortURL, 0, len(requests))
	for _, item := range requests {
		id, err := s.urlService.GenerateShortURL(item.OriginalURL, userID)
		if err != nil {
			return nil, err
		}

		resp = append(resp, models.BatchShortURL{
			CorrelationID: item.CorrelationID,
			ShortURL:      s.baseURL + "/" + id,
		})
	}

	return resp, nil
}

// GetUserURLs получает все URL пользователя
func (s *CommonURLService) GetUserURLs(userID string) ([]models.URLEntry, error) {
	if userID == "" {
		return nil, errors.New("userID не может быть пустым")
	}

	urls, err := s.urlService.GetURLsByUser(userID)
	if err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		return []models.URLEntry{}, nil
	}

	resp := make([]models.URLEntry, 0, len(urls))
	for _, u := range urls {
		resp = append(resp, models.URLEntry{
			ShortURL:    s.baseURL + "/" + u.ShortURL,
			OriginalURL: u.OriginalURL,
		})
	}

	return resp, nil
}

// BatchDeleteURLs помечает URL пользователя на удаление асинхронно
func (s *CommonURLService) BatchDeleteURLs(userID string, ids []string) error {
	if userID == "" {
		return errors.New("userID не может быть пустым")
	}
	if len(ids) == 0 {
		return errors.New("список ID не может быть пустым")
	}

	if s.deleteChan == nil {
		return errors.New("канал удаления не настроен")
	}

	task := async.DeleteTask{
		UserID:    userID,
		ShortURLs: ids,
	}

	select {
	case s.deleteChan <- task:
		return nil
	default:
		return errors.New("не удалось поставить задачу на удаление в очередь")
	}
}

// GetStats получает статистику сервиса
func (s *CommonURLService) GetStats() (urlCount int, userCount int, err error) {
	return s.urlService.GetStats()
}