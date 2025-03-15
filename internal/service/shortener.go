package service

import (
	"context"

	"github.com/Levan1e/url-shortener-service/internal/domain"
	"github.com/Levan1e/url-shortener-service/internal/utils"
)

const (
	shortURLLength = 10
	maxAttempts    = 5
)

type Storage interface {
	Save(ctx context.Context, originalURL, shortURL string) error
	GetShort(ctx context.Context, originalURL string) (string, error)
	GetOriginal(ctx context.Context, shortURL string) (string, error)
}

type ShortenerService struct {
	storage Storage
}

func NewShortenerService(storage Storage) *ShortenerService {
	return &ShortenerService{
		storage: storage,
	}
}

func (s *ShortenerService) GetShortenByOriginal(ctx context.Context, originalURL string) (string, error) {
	if short, err := s.storage.GetShort(ctx, originalURL); err != nil {
		return "", err
	} else if short != "" {
		return short, nil
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		shortUrl, err := utils.GenerateRandomString(shortURLLength)
		if err != nil {
			return "", err
		}
		if err := s.storage.Save(ctx, originalURL, shortUrl); err != nil {
			if err == domain.ErrAlreadyExist {
				continue
			}
			return "", err
		}
		return shortUrl, nil
	}
	return "", domain.InternalServerError
}

func (s *ShortenerService) GetOriginalByShorten(ctx context.Context, shortURL string) (string, error) {
	original, err := s.storage.GetOriginal(ctx, shortURL)
	if err != nil {
		return "", err
	}
	if original == "" {
		return "", domain.UrlNotFound
	}
	return original, nil
}
