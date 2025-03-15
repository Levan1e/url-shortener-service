package postgres

import (
	"context"

	"github.com/Levan1e/url-shortener-service/internal/domain"
	postgres_helpers "github.com/Levan1e/url-shortener-service/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *PostgresStorage {
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
	return postgres_helpers.Select[string](ctx, s.pool, query, originalURL)
}

func (s *PostgresStorage) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	query := `SELECT original_url FROM urls WHERE short_url = $1;`
	return postgres_helpers.Select[string](ctx, s.pool, query, shortURL)
}