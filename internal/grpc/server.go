// Package grpc содержит gRPC сервер и реализацию URL shortener сервиса.
package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

// URLShortenerServer реализует gRPC сервис для сокращения URL
type URLShortenerServer struct {
	UnimplementedURLShortenerServiceServer
	commonService *services.CommonURLService
}

// NewURLShortenerServer создает новый gRPC сервер
func NewURLShortenerServer(commonService *services.CommonURLService) *URLShortenerServer {
	return &URLShortenerServer{
		commonService: commonService,
	}
}

// CreateShortURL создает короткую ссылку из обычного URL
func (s *URLShortenerServer) CreateShortURL(ctx context.Context, req *CreateShortURLRequest) (*CreateShortURLResponse, error) {
	if req.URL == "" {
		return nil, status.Error(codes.InvalidArgument, "URL не может быть пустым")
	}

	shortURL, alreadyExists, err := s.commonService.CreateShortURL(req.URL, req.UserID)
	if err != nil {
		log.Printf("Ошибка создания короткой ссылки: %v", err)
		return nil, status.Error(codes.Internal, "Ошибка создания короткой ссылки")
	}

	return &CreateShortURLResponse{
		ShortURL:      shortURL,
		AlreadyExists: alreadyExists,
	}, nil
}

// CreateShortJSONURL создает короткую ссылку в формате JSON
func (s *URLShortenerServer) CreateShortJSONURL(ctx context.Context, req *CreateShortJSONURLRequest) (*CreateShortJSONURLResponse, error) {
	if req.URL == "" {
		return nil, status.Error(codes.InvalidArgument, "URL не может быть пустым")
	}

	shortURL, alreadyExists, err := s.commonService.CreateShortURL(req.URL, req.UserID)
	if err != nil {
		log.Printf("Ошибка создания короткой ссылки: %v", err)
		return nil, status.Error(codes.Internal, "Ошибка создания короткой ссылки")
	}

	return &CreateShortJSONURLResponse{
		Result:        shortURL,
		AlreadyExists: alreadyExists,
	}, nil
}

// GetShortURL получает оригинальный URL по короткому ID
func (s *URLShortenerServer) GetShortURL(ctx context.Context, req *GetShortURLRequest) (*GetShortURLResponse, error) {
	if req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID не может быть пустым")
	}

	originalURL, deleted, err := s.commonService.GetOriginalURL(req.ID)
	if err != nil {
		log.Printf("Ошибка получения URL: %v", err)
		return nil, status.Error(codes.NotFound, "URL не найден")
	}

	if deleted {
		return nil, status.Error(codes.FailedPrecondition, "URL был удален")
	}

	return &GetShortURLResponse{
		OriginalURL: originalURL,
		Deleted:     deleted,
	}, nil
}

// Ping проверяет подключение к базе данных
func (s *URLShortenerServer) Ping(ctx context.Context, req *PingRequest) (*PingResponse, error) {
	err := s.commonService.Ping()
	if err != nil {
		log.Printf("Ошибка проверки соединения: %v", err)
		return &PingResponse{Success: false}, nil
	}

	return &PingResponse{Success: true}, nil
}

// Batch создает несколько коротких ссылок за один запрос
func (s *URLShortenerServer) Batch(ctx context.Context, req *BatchRequest) (*BatchResponse, error) {
	if len(req.Urls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Список URL не может быть пустым")
	}

	// Преобразуем gRPC структуры в модели бизнес-логики
	modelRequests := make([]models.BatchOriginalURL, len(req.Urls))
	for i, url := range req.Urls {
		modelRequests[i] = models.BatchOriginalURL{
			CorrelationID: url.CorrelationID,
			OriginalURL:   url.OriginalURL,
		}
	}

	results, err := s.commonService.BatchCreateShortURLs(modelRequests, req.UserID)
	if err != nil {
		log.Printf("Ошибка пакетного создания URL: %v", err)
		return nil, status.Error(codes.Internal, "Ошибка пакетного создания URL")
	}

	// Преобразуем результаты обратно в gRPC структуры
	grpcResults := make([]*BatchShortURL, len(results))
	for i, result := range results {
		grpcResults[i] = &BatchShortURL{
			CorrelationID: result.CorrelationID,
			ShortURL:      result.ShortURL,
		}
	}

	return &BatchResponse{Urls: grpcResults}, nil
}

// GetUserURLs получает все URL пользователя
func (s *URLShortenerServer) GetUserURLs(ctx context.Context, req *GetUserURLsRequest) (*GetUserURLsResponse, error) {
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID не может быть пустым")
	}

	urls, err := s.commonService.GetUserURLs(req.UserID)
	if err != nil {
		log.Printf("Ошибка получения URL пользователя: %v", err)
		return nil, status.Error(codes.Internal, "Ошибка получения URL пользователя")
	}

	// Преобразуем в gRPC структуры
	grpcURLs := make([]*URLEntry, len(urls))
	for i, url := range urls {
		grpcURLs[i] = &URLEntry{
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
		}
	}

	return &GetUserURLsResponse{Urls: grpcURLs}, nil
}

// BatchDelete удаляет несколько URL пользователя
func (s *URLShortenerServer) BatchDelete(ctx context.Context, req *BatchDeleteRequest) (*BatchDeleteResponse, error) {
	if req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID не может быть пустым")
	}

	if len(req.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Список ID не может быть пустым")
	}

	err := s.commonService.BatchDeleteURLs(req.UserID, req.Ids)
	if err != nil {
		log.Printf("Ошибка пакетного удаления URL: %v", err)
		// Если канал удаления не настроен, возвращаем ошибку недоступности
		if errors.Is(err, errors.New("канал удаления не настроен")) {
			return nil, status.Error(codes.Unavailable, "Сервис удаления недоступен")
		}
		return nil, status.Error(codes.Internal, "Ошибка удаления URL")
	}

	return &BatchDeleteResponse{Success: true}, nil
}

// GetStats получает статистику сервиса
func (s *URLShortenerServer) GetStats(ctx context.Context, req *GetStatsRequest) (*GetStatsResponse, error) {
	urlCount, userCount, err := s.commonService.GetStats()
	if err != nil {
		log.Printf("Ошибка получения статистики: %v", err)
		return nil, status.Error(codes.Internal, "Ошибка получения статистики")
	}

	return &GetStatsResponse{
		Urls:  int32(urlCount),
		Users: int32(userCount),
	}, nil
}