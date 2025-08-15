package models

// Переменные модели URL
type (
	ShortenRequest struct {
		URL string `json:"url"`
	}
	ShortenResponse struct {
		Result string `json:"result"`
	}
	BatchOriginalURL struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
	BatchShortURL struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)
