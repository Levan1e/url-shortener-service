package memory_test

import (
	"context"
	"sync"
	"testing"

	"github.com/Levan1e/url-shortener-service/internal/domain"
	"github.com/Levan1e/url-shortener-service/internal/repository/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_Save_Success(t *testing.T) {
	storage := memory.NewStorage()
	err := storage.Save(context.Background(), "http://example.com", "abc123")
	assert.NoError(t, err)

	short, err := storage.GetShort(context.Background(), "http://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", short)
}

func TestMemoryStorage_Save_Duplicate(t *testing.T) {
	storage := memory.NewStorage()
	err := storage.Save(context.Background(), "http://example.com", "abc123")
	assert.NoError(t, err)

	err = storage.Save(context.Background(), "http://example.com", "abc123")
	assert.ErrorIs(t, err, domain.ErrAlreadyExist)
}

func TestMemoryStorage_GetShort_NotFound(t *testing.T) {
	storage := memory.NewStorage()
	short, err := storage.GetShort(context.Background(), "http://unknown.com")
	assert.NoError(t, err)
	assert.Empty(t, short)
}

func TestMemoryStorage_GetOriginal_NotFound(t *testing.T) {
	storage := memory.NewStorage()
	original, err := storage.GetOriginal(context.Background(), "xyz123")

	assert.NoError(t, err)
	assert.Empty(t, original)
}

func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	storage := memory.NewStorage()
	const goroutines = 10
	var wg sync.WaitGroup

	saveFunc := func(id int) {
		defer wg.Done()
		err := storage.Save(context.Background(), "http://example.com", "abc123")
		if err != nil && err != domain.ErrAlreadyExist {
			t.Errorf("unexpected error: %v", err)
		}
	}

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go saveFunc(i)
	}
	wg.Wait()

	short, err := storage.GetShort(context.Background(), "http://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", short)
}
