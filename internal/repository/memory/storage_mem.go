package memory

import (
	"context"
	"sync"

	"github.com/Levan1e/url-shortener-service/internal/domain"
)

type MemoryStorage struct {
	mu              sync.RWMutex
	originalToShort map[string]string
	shortToOriginal map[string]string
}

func NewStorage() *MemoryStorage {
	return &MemoryStorage{
		originalToShort: make(map[string]string),
		shortToOriginal: make(map[string]string),
	}
}

func (s *MemoryStorage) Save(_ context.Context, originalURL, shortURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exist := s.shortToOriginal[shortURL]; exist {
		return domain.ErrAlreadyExist
	}
	s.originalToShort[originalURL] = shortURL
	s.shortToOriginal[shortURL] = originalURL
	return nil
}

func (s *MemoryStorage) GetShort(_ context.Context, originalURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if short, exists := s.originalToShort[originalURL]; exists {
		return short, nil
	}
	return "", nil
}

func (s *MemoryStorage) GetOriginal(_ context.Context, shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if original, exists := s.shortToOriginal[shortURL]; exists {
		return original, nil
	}
	return "", nil
}
