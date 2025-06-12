package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
)

type DBStore struct {
	db *pgxpool.Pool
}

func NewDBStore(db *pgxpool.Pool) *DBStore {
	return &DBStore{db: db}
}

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

func (d *DBStore) Get(shortID string) (string, error) {
	const query = `SELECT original_url FROM short_urls WHERE short_id = $1;`
	var url string
	err := d.db.QueryRow(context.Background(), query, shortID).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (d *DBStore) SaveToFile(path string) error {
	return nil
}

func (d *DBStore) LoadFromFile(path string) error {
	return nil
}

func (d *DBStore) Ping() error {
	return d.db.Ping(context.Background())
}

func (d *DBStore) GetShortIDByOriginalURL(originalURL string) (string, error) {
	const query = `SELECT short_id FROM short_urls WHERE original_url = $1`
	var shortID string
	err := d.db.QueryRow(context.Background(), query, originalURL).Scan(&shortID)
	if err != nil {
		return "", err
	}
	return shortID, nil
}

func (d *DBStore) GetByUser(userID string) ([]models.URLEntry, error) {
	const query = `SELECT short_id, original_url FROM short_urls WHERE user_uuid = $1;`

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
