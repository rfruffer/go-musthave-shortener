package posgreconfig

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing database config: %w", err)
	}
	config.MaxConns = 10
	config.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	const createTable = `
	CREATE TABLE IF NOT EXISTS short_urls (
		short_id TEXT PRIMARY KEY,
		original_url TEXT NOT NULL,
		user_uuid UUID NOT NULL
	);
	`
	_, err = pool.Exec(context.Background(), createTable)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Database connection established")
	return pool, nil
}

func CloseDB(pool *pgxpool.Pool) {
	if pool == nil {
		log.Println("No database connection to close")
		return
	}

	pool.Close()
	log.Println("Database connection closed")
}
