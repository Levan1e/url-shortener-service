package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Levan1e/url-shortener-service/internal/repository/postgres"
)

func getTestStorage(t *testing.T) *postgres.PostgresStorage {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		t.Skip("DATABASE_URL не установлен; пропускаем интеграционные тесты для PostgreSQL")
	}

	storage, err := postgres.NewStorage(connStr)
	if err != nil {
		t.Fatalf("не удалось создать PostgresStorage: %v", err)
	}

	cleanupTable(storage, t)
	return storage
}

func cleanupTable(storage *postgres.PostgresStorage, t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := storage.Pool().Exec(ctx, "TRUNCATE TABLE urls;")
	if err != nil {
		t.Fatalf("не удалось очистить таблицу urls: %v", err)
	}
}

func TestPostgresStorage_SaveAndRetrieve(t *testing.T) {
	storage := getTestStorage(t)
	original := "https://example.com"
	short := "AbCdEfGh12"

	err := storage.Save(original, short)
	assert.NoError(t, err, "Save должен выполниться без ошибки")

	gotShort, err := storage.GetShort(original)
	assert.NoError(t, err)
	assert.Equal(t, short, gotShort)

	gotOriginal, err := storage.GetOriginal(short)
	assert.NoError(t, err)
	assert.Equal(t, original, gotOriginal)
}

func TestPostgresStorage_DuplicateSave(t *testing.T) {
	storage := getTestStorage(t)
	original := "https://example.com"
	short := "AbCdEfGh12"

	err := storage.Save(original, short)
	assert.NoError(t, err)

	err = storage.Save(original, short)
	assert.NoError(t, err)
}

func TestPostgresStorage_ErrorOnMismatch(t *testing.T) {
	storage := getTestStorage(t)
	original := "https://example.com"
	short1 := "AbCdEfGh12"
	short2 := "ZyXwVuTsRq"

	err := storage.Save(original, short1)
	assert.NoError(t, err)

	err = storage.Save(original, short2)
	assert.Error(t, err)
	assert.Equal(t, "оригинальный URL уже сопоставлен с другим коротким URL", err.Error())
}

func TestPostgresStorage_GetNotFound(t *testing.T) {
	storage := getTestStorage(t)

	_, err := storage.GetShort("https://nonexistent.com")
	assert.Error(t, err)

	_, err = storage.GetOriginal("NonExistent")
	assert.Error(t, err)
}
