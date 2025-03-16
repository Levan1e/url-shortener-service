package memory

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/Levan1e/url-shortener-service/internal/domain"
)

type MemoryStorage struct {
	mu              sync.RWMutex
	originalToShort map[string]string
	shortToOriginal map[string]string
}

func NewStorage() *MemoryStorage {
	storage := &MemoryStorage{
		originalToShort: make(map[string]string),
		shortToOriginal: make(map[string]string),
	}
	storage.loadFromFile("storage.json")
	return storage
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

func (s *MemoryStorage) loadFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	type snapshot struct {
		OriginalToShort map[string]string `json:"original_to_short"`
		ShortToOriginal map[string]string `json:"short_to_original"`
	}

	var snap snapshot
	if err := json.NewDecoder(file).Decode(&snap); err == nil {
		s.originalToShort = snap.OriginalToShort
		s.shortToOriginal = snap.ShortToOriginal
	}
}

func (s *MemoryStorage) SaveToFileOnShutdown(filename string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type snapshot struct {
		OriginalToShort map[string]string `json:"original_to_short"`
		ShortToOriginal map[string]string `json:"short_to_original"`
	}

	snap := snapshot{
		OriginalToShort: s.originalToShort,
		ShortToOriginal: s.shortToOriginal,
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(&snap)
}
