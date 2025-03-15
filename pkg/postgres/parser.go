package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Select[T any](ctx context.Context, db *pgxpool.Pool, query string, args ...any) (T, error) {
	var target T
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return target, err
	}
	target, err = pgx.CollectOneRow(rows, pgx.RowTo[T])
	if err == pgx.ErrNoRows {
		return target, nil
	}
	return target, err
}

func SelectMany[T any](ctx context.Context, db *pgxpool.Pool, query string, args ...any) ([]T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	target, err := pgx.CollectRows(rows, pgx.RowTo[T])
	if err == pgx.ErrNoRows {
		return target, nil
	}
	return target, err
}

func SelectStruct[T any](ctx context.Context, db *pgxpool.Pool, query string, args ...any) (T, error) {
	var target T
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return target, err
	}
	target, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err == pgx.ErrNoRows {
		return target, nil
	}
	return target, err
}

func SelectStructMany[T any](ctx context.Context, db *pgxpool.Pool, query string, args ...any) ([]T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	target, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err == pgx.ErrNoRows {
		return target, nil
	}
	return target, err
}
