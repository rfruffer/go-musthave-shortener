// Package repository предоставляет абстракции и реализации хранилищ (Postgres, файл).
package repository

import (
	"context"

	// "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
)

// DBStore инициализирует пакет pgxpool
type DBStore struct {
	db *pgxpool.Pool
}

// NewDBStore создает новый DBStore
func NewDBStore(db *pgxpool.Pool) *DBStore {
	return &DBStore{db: db}
}

// Save функция сохранения в базу
func (d *DBStore) Save(shortID, originalURL, uuid string) error {
	const query = `
		INSERT INTO short_urls (short_id, original_url, user_uuid)
		VALUES ($1, $2, $3);
	`
	_, err := d.db.Exec(context.Background(), query, shortID, originalURL, uuid)
	if err != nil {
		return err
	}
	return nil
}

// GetURLByShort функция получения URL
func (d *DBStore) GetURLByShort(shortID string) (models.URLEntry, error) {
	const query = `SELECT user_uuid, short_id, original_url, is_deleted FROM short_urls WHERE short_id = $1;`
	var entry models.URLEntry
	err := d.db.QueryRow(context.Background(), query, shortID).Scan(
		&entry.UUID,
		&entry.ShortURL,
		&entry.OriginalURL,
		&entry.DeletedFlag,
	)
	if err != nil {
		return models.URLEntry{}, err
	}

	return entry, nil
}

// SaveToFile сохранение в файл
func (d *DBStore) SaveToFile(path string) error {
	return nil
}

// LoadFromFile загрузка из файла
func (d *DBStore) LoadFromFile(path string) error {
	return nil
}

// Ping проверка подключения
func (d *DBStore) Ping() error {
	return d.db.Ping(context.Background())
}

// GetShortIDByOriginalURL получить коротку ссылку по URL
func (d *DBStore) GetShortIDByOriginalURL(originalURL string) (string, error) {
	const query = `SELECT short_id FROM short_urls WHERE original_url = $1`
	var shortID string
	err := d.db.QueryRow(context.Background(), query, originalURL).Scan(&shortID)
	if err != nil {
		return "", err
	}
	return shortID, nil
}

// GetByUser получить юзера
func (d *DBStore) GetByUser(userID string) ([]models.URLEntry, error) {
	const query = `SELECT short_id, original_url 
	FROM short_urls 
	WHERE user_uuid = $1 and is_deleted IS NOT TRUE;`

	rows, err := d.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.URLEntry
	for rows.Next() {
		var url models.URLEntry
		if err := rows.Scan(&url.ShortURL, &url.OriginalURL); err != nil {
			return nil, err
		}
		results = append(results, url)
	}
	return results, nil
}

// MarkURLsDeleted пометить как удаленные
func (d *DBStore) MarkURLsDeleted(userID string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	query := `UPDATE short_urls 
	SET is_deleted = true 
	WHERE user_uuid = $1 AND short_id = ANY($2);`
	_, err := d.db.Exec(context.Background(), query, userID, ids)
	return err
}
