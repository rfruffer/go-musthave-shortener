// Package grpc содержит gRPC сервис интерфейсы
package grpc

import (
	"context"
	"errors"
)

var ErrMethodNotImplemented = errors.New("method not implemented")

// URLShortenerServiceServer интерфейс для gRPC сервиса
type URLShortenerServiceServer interface {
	CreateShortURL(context.Context, *CreateShortURLRequest) (*CreateShortURLResponse, error)
	CreateShortJSONURL(context.Context, *CreateShortJSONURLRequest) (*CreateShortJSONURLResponse, error)
	GetShortURL(context.Context, *GetShortURLRequest) (*GetShortURLResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Batch(context.Context, *BatchRequest) (*BatchResponse, error)
	GetUserURLs(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error)
	BatchDelete(context.Context, *BatchDeleteRequest) (*BatchDeleteResponse, error)
	GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error)
}

// URLShortenerServiceClient интерфейс для gRPC клиента
type URLShortenerServiceClient interface {
	CreateShortURL(ctx context.Context, req *CreateShortURLRequest) (*CreateShortURLResponse, error)
	CreateShortJSONURL(ctx context.Context, req *CreateShortJSONURLRequest) (*CreateShortJSONURLResponse, error)
	GetShortURL(ctx context.Context, req *GetShortURLRequest) (*GetShortURLResponse, error)
	Ping(ctx context.Context, req *PingRequest) (*PingResponse, error)
	Batch(ctx context.Context, req *BatchRequest) (*BatchResponse, error)
	GetUserURLs(ctx context.Context, req *GetUserURLsRequest) (*GetUserURLsResponse, error)
	BatchDelete(ctx context.Context, req *BatchDeleteRequest) (*BatchDeleteResponse, error)
	GetStats(ctx context.Context, req *GetStatsRequest) (*GetStatsResponse, error)
}

// UnimplementedURLShortenerServiceServer базовая реализация для обратной совместимости
type UnimplementedURLShortenerServiceServer struct{}

func (UnimplementedURLShortenerServiceServer) CreateShortURL(context.Context, *CreateShortURLRequest) (*CreateShortURLResponse, error) {
	return nil, ErrMethodNotImplemented
}
func (UnimplementedURLShortenerServiceServer) CreateShortJSONURL(context.Context, *CreateShortJSONURLRequest) (*CreateShortJSONURLResponse, error) {
	return nil, ErrMethodNotImplemented
}
func (UnimplementedURLShortenerServiceServer) GetShortURL(context.Context, *GetShortURLRequest) (*GetShortURLResponse, error) {
	return nil, ErrMethodNotImplemented
}
func (UnimplementedURLShortenerServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, ErrMethodNotImplemented
}
func (UnimplementedURLShortenerServiceServer) Batch(context.Context, *BatchRequest) (*BatchResponse, error) {
	return nil, ErrMethodNotImplemented
}
func (UnimplementedURLShortenerServiceServer) GetUserURLs(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error) {
	return nil, ErrMethodNotImplemented
}
func (UnimplementedURLShortenerServiceServer) BatchDelete(context.Context, *BatchDeleteRequest) (*BatchDeleteResponse, error) {
	return nil, ErrMethodNotImplemented
}
func (UnimplementedURLShortenerServiceServer) GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error) {
	return nil, ErrMethodNotImplemented
}

// Фиктивная функция регистрации для совместимости
func RegisterURLShortenerServiceServer(s interface{}, srv URLShortenerServiceServer) {
}

// Фиктивная функция создания клиента для совместимости
func NewURLShortenerServiceClient(conn interface{}) URLShortenerServiceClient {
	return nil
}