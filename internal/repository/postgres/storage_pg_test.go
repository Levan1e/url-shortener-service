package postgres_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"github.com/Levan1e/url-shortener-service/internal/domain"
	"github.com/Levan1e/url-shortener-service/internal/repository/postgres"
	"github.com/Levan1e/url-shortener-service/internal/repository/postgres/mocks"
)

type MockRow struct {
	scanFunc func(dest ...interface{}) error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	return m.scanFunc(dest...)
}

func NewMockRow(scanFunc func(dest ...interface{}) error) pgx.Row {
	return &MockRow{scanFunc: scanFunc}
}

func TestPostgresStorage_Save_Duplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := mocks.NewMockPgxPoolInterface(ctrl)
	storage := postgres.NewStorage(mockPool)

	query := `
	INSERT INTO urls (original_url, short_url)
	VALUES ($1, $2)
	ON CONFLICT (short_url) DO NOTHING;
	`

	mockPool.EXPECT().
		Exec(gomock.Any(), gomock.Eq(query), gomock.Eq("http://example.com"), gomock.Eq("abc123")).
		Return(pgconn.NewCommandTag("INSERT 0"), nil)

	err := storage.Save(context.Background(), "http://example.com", "abc123")
	assert.ErrorIs(t, err, domain.ErrAlreadyExist)
}

func TestPostgresStorage_GetShort_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := mocks.NewMockPgxPoolInterface(ctrl)
	storage := postgres.NewStorage(mockPool)

	query := `SELECT short_url FROM urls WHERE original_url = $1;`

	mockRow := NewMockRow(func(dest ...interface{}) error {
		*dest[0].(*string) = "abc123"
		return nil
	})

	mockPool.EXPECT().
		QueryRow(context.Background(), query, "http://example.com").
		Return(mockRow)

	short, err := storage.GetShort(context.Background(), "http://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", short)
}

func TestPostgresStorage_GetShort_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := mocks.NewMockPgxPoolInterface(ctrl)
	storage := postgres.NewStorage(mockPool)

	query := `SELECT short_url FROM urls WHERE original_url = $1;`
	mockRow := NewMockRow(func(dest ...interface{}) error {
		return pgx.ErrNoRows
	})

	mockPool.EXPECT().
		QueryRow(context.Background(), query, "http://unknown.com").
		Return(mockRow)

	short, err := storage.GetShort(context.Background(), "http://unknown.com")
	assert.NoError(t, err)
	assert.Empty(t, short)
}

func TestPostgresStorage_GetOriginal_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := mocks.NewMockPgxPoolInterface(ctrl)
	storage := postgres.NewStorage(mockPool)

	query := `SELECT original_url FROM urls WHERE short_url = $1;`

	mockRow := NewMockRow(func(dest ...interface{}) error {
		*dest[0].(*string) = "http://example.com"
		return nil
	})

	mockPool.EXPECT().
		QueryRow(context.Background(), query, "abc123").
		Return(mockRow)

	original, err := storage.GetOriginal(context.Background(), "abc123")
	assert.NoError(t, err)
	assert.Equal(t, "http://example.com", original)
}

func TestPostgresStorage_GetOriginal_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := mocks.NewMockPgxPoolInterface(ctrl)
	storage := postgres.NewStorage(mockPool)

	query := `SELECT original_url FROM urls WHERE short_url = $1;`
	mockRow := NewMockRow(func(dest ...interface{}) error {
		return pgx.ErrNoRows
	})

	mockPool.EXPECT().
		QueryRow(context.Background(), query, "xyz123").
		Return(mockRow)

	original, err := storage.GetOriginal(context.Background(), "xyz123")
	assert.NoError(t, err)
	assert.Empty(t, original)
}
