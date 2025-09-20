// Package grpc содержит типы данных для gRPC API
package grpc

// CreateShortURLRequest запрос для создания короткой ссылки
type CreateShortURLRequest struct {
	Url    string `json:"url"`
	UserId string `json:"user_id"`
}

// CreateShortURLResponse ответ для создания короткой ссылки
type CreateShortURLResponse struct {
	ShortUrl      string `json:"short_url"`
	AlreadyExists bool   `json:"already_exists"`
}

// CreateShortJSONURLRequest запрос для создания короткой ссылки в JSON формате
type CreateShortJSONURLRequest struct {
	Url    string `json:"url"`
	UserId string `json:"user_id"`
}

// CreateShortJSONURLResponse ответ для создания короткой ссылки в JSON формате
type CreateShortJSONURLResponse struct {
	Result        string `json:"result"`
	AlreadyExists bool   `json:"already_exists"`
}

// GetShortURLRequest запрос для получения оригинального URL
type GetShortURLRequest struct {
	Id string `json:"id"`
}

// GetShortURLResponse ответ для получения оригинального URL
type GetShortURLResponse struct {
	OriginalUrl string `json:"original_url"`
	Deleted     bool   `json:"deleted"`
}

// PingRequest запрос для проверки подключения
type PingRequest struct{}

// PingResponse ответ для проверки подключения
type PingResponse struct {
	Success bool `json:"success"`
}

// BatchOriginalURL элемент пакетного запроса
type BatchOriginalURL struct {
	CorrelationId string `json:"correlation_id"`
	OriginalUrl   string `json:"original_url"`
}

// BatchShortURL элемент пакетного ответа
type BatchShortURL struct {
	CorrelationId string `json:"correlation_id"`
	ShortUrl      string `json:"short_url"`
}

// BatchRequest запрос для пакетного создания
type BatchRequest struct {
	Urls   []*BatchOriginalURL `json:"urls"`
	UserId string              `json:"user_id"`
}

// BatchResponse ответ для пакетного создания
type BatchResponse struct {
	Urls []*BatchShortURL `json:"urls"`
}

// GetUserURLsRequest запрос для получения URL пользователя
type GetUserURLsRequest struct {
	UserId string `json:"user_id"`
}

// URLEntry URL пользователя
type URLEntry struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

// GetUserURLsResponse ответ для получения URL пользователя
type GetUserURLsResponse struct {
	Urls []*URLEntry `json:"urls"`
}

// BatchDeleteRequest запрос для пакетного удаления
type BatchDeleteRequest struct {
	Ids    []string `json:"ids"`
	UserId string   `json:"user_id"`
}

// BatchDeleteResponse ответ для пакетного удаления
type BatchDeleteResponse struct {
	Success bool `json:"success"`
}

// GetStatsRequest запрос для получения статистики
type GetStatsRequest struct{}

// GetStatsResponse ответ для получения статистики
type GetStatsResponse struct {
	Urls  int32 `json:"urls"`
	Users int32 `json:"users"`
}