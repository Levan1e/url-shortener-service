package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Levan1e/url-shortener-service/internal/repository/mocks"
	"github.com/Levan1e/url-shortener-service/internal/service"
)

func TestShorten_ReturnsExistingShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	originalURL := "https://example.com"
	existingShort := "AbCdEfGh12"

	mockStorage.
		EXPECT().
		GetShort(originalURL).
		Return(existingShort, nil).
		Times(1)

	svc := service.NewShortenerService(mockStorage)
	result, err := svc.Shorten(originalURL)
	assert.NoError(t, err)
	assert.Equal(t, existingShort, result)
}

func TestShorten_GeneratesNewShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	originalURL := "https://example.com"
	mockStorage := mocks.NewMockStorage(ctrl)

	mockStorage.
		EXPECT().
		GetShort(originalURL).
		Return("", errors.New("not found")).
		Times(1)

	mockStorage.
		EXPECT().
		Save(originalURL, gomock.Any()).
		DoAndReturn(func(url, short string) error {
			if len(short) == 10 {
				return nil
			}
			return errors.New("invalid short url length")
		}).
		Times(1)

	svc := service.NewShortenerService(mockStorage)
	result, err := svc.Shorten(originalURL)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(result))
}

func TestShorten_FailsAfterMaxAttempts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	originalURL := "https://example.com"
	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.
		EXPECT().
		GetShort(originalURL).
		Return("", errors.New("not found")).
		Times(1)

	mockStorage.
		EXPECT().
		Save(originalURL, gomock.Any()).
		Return(errors.New("collision")).
		Times(5)

	svc := service.NewShortenerService(mockStorage)
	result, err := svc.Shorten(originalURL)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, "не удалось сгенерировать уникальный короткий URL", err.Error())
}

func TestResolve_ReturnsOriginalURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shortURL := "AbCdEfGh12"
	originalURL := "https://example.com"
	mockStorage := mocks.NewMockStorage(ctrl)

	mockStorage.
		EXPECT().
		GetOriginal(shortURL).
		Return(originalURL, nil).
		Times(1)

	svc := service.NewShortenerService(mockStorage)
	result, err := svc.Resolve(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, result)
}
