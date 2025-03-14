package memory

import (
	"errors"
	"sync"
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

func (s *MemoryStorage) Save(originalURL, shortURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, exists := s.shortToOriginal[shortURL]; exists {
		if existing != originalURL {
			return errors.New("короткий URL уже используется для другого оригинального URL")
		}
		return nil
	}

	if existing, exists := s.originalToShort[originalURL]; exists {
		if existing != shortURL {
			return errors.New("оригинальный URL уже сопоставлен с другим коротким URL")
		}
		return nil
	}

	s.originalToShort[originalURL] = shortURL
	s.shortToOriginal[shortURL] = originalURL

	return nil
}

func (s *MemoryStorage) GetShort(originalURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if short, exists := s.originalToShort[originalURL]; exists {
		return short, nil
	}
	return "", errors.New("короткий URL для данного оригинального URL не найден")
}

func (s *MemoryStorage) GetOriginal(shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if original, exists := s.shortToOriginal[shortURL]; exists {
		return original, nil
	}
	return "", errors.New("оригинальный URL для данного короткого URL не найден")
}
