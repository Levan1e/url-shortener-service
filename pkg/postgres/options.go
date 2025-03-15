package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func NewPostgresPool(ctx context.Context, config *Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.Database)
	return pgxpool.New(ctx, connString)
}

func Migrate(pool *pgxpool.Pool, migrationsDir string) error {
	return goose.Up(stdlib.OpenDBFromPool(pool), migrationsDir)
}
