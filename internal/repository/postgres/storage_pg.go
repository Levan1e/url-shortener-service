package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewStorage(connString string) (*PostgresStorage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := createTable(context.Background(), pool); err != nil {
		pool.Close()
		return nil, err
	}

	return &PostgresStorage{pool: pool}, nil
}

func createTable(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		original_url TEXT PRIMARY KEY,
		short_url TEXT UNIQUE NOT NULL
	);
	`
	_, err := pool.Exec(ctx, query)
	return err
}

func (s *PostgresStorage) Save(originalURL, shortURL string) error {
	ctx := context.Background()
	query := `
	INSERT INTO urls (original_url, short_url)
	VALUES ($1, $2)
	ON CONFLICT (original_url) DO NOTHING;
	`
	ct, err := s.pool.Exec(ctx, query, originalURL, shortURL)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		existing, err := s.GetShort(originalURL)
		if err != nil {
			return err
		}
		if existing != shortURL {
			return errors.New("оригинальный URL уже сопоставлен с другим коротким URL")
		}
	}
	return nil
}

func (s *PostgresStorage) GetShort(originalURL string) (string, error) {
	ctx := context.Background()
	query := `SELECT short_url FROM urls WHERE original_url = $1;`
	var shortURL string
	err := s.pool.QueryRow(ctx, query, originalURL).Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (s *PostgresStorage) GetOriginal(shortURL string) (string, error) {
	ctx := context.Background()
	query := `SELECT original_url FROM urls WHERE short_url = $1;`
	var originalURL string
	err := s.pool.QueryRow(ctx, query, shortURL).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (s *PostgresStorage) Pool() *pgxpool.Pool {
	return s.pool
}
