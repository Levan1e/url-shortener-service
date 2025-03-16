package postgres

import (
	"context"

	"github.com/Levan1e/url-shortener-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxPoolInterface interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type PostgresStorage struct {
	pool PgxPoolInterface
}

func NewStorage(pool PgxPoolInterface) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (s *PostgresStorage) Save(ctx context.Context, originalURL, shortURL string) error {
	query := `
	INSERT INTO urls (original_url, short_url)
	VALUES ($1, $2)
	ON CONFLICT (short_url) DO NOTHING;
	`
	res, err := s.pool.Exec(ctx, query, originalURL, shortURL)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return domain.ErrAlreadyExist
	}
	return nil
}

func (s *PostgresStorage) GetShort(ctx context.Context, originalURL string) (string, error) {
	query := `SELECT short_url FROM urls WHERE original_url = $1;`
	var shortURL string
	err := s.pool.QueryRow(ctx, query, originalURL).Scan(&shortURL)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return shortURL, nil
}

func (s *PostgresStorage) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	query := `SELECT original_url FROM urls WHERE short_url = $1;`
	var originalURL string
	err := s.pool.QueryRow(ctx, query, shortURL).Scan(&originalURL)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return originalURL, nil
}
