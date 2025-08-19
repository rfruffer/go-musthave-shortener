// Package models содержит модели данных для сервиса сокращения URL.
package models

// Переменные модели URL
type URLEntry struct {
	UUID        string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"original_url"`
	DeletedFlag bool   `db:"is_deleted"`
}
