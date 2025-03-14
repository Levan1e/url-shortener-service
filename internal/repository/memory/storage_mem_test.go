package memory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Levan1e/url-shortener-service/internal/repository/memory"
)

func TestMemoryStorage_SaveAndRetrieve(t *testing.T) {
	store := memory.NewStorage()
	original := "https://example.com"
	short := "AbCdEfGh12"

	err := store.Save(original, short)
	assert.NoError(t, err)

	retrievedShort, err := store.GetShort(original)
	assert.NoError(t, err)
	assert.Equal(t, short, retrievedShort)

	retrievedOriginal, err := store.GetOriginal(short)
	assert.NoError(t, err)
	assert.Equal(t, original, retrievedOriginal)
}

func TestMemoryStorage_DuplicateSave(t *testing.T) {
	store := memory.NewStorage()
	original := "https://example.com"
	short := "AbCdEfGh12"

	err := store.Save(original, short)
	assert.NoError(t, err)

	err = store.Save(original, short)
	assert.NoError(t, err)
}

func TestMemoryStorage_ErrorOnMismatch(t *testing.T) {
	store := memory.NewStorage()
	original := "https://example.com"
	short1 := "AbCdEfGh12"
	short2 := "ZyXwVuTsRq"

	err := store.Save(original, short1)
	assert.NoError(t, err)

	err = store.Save(original, short2)
	assert.Error(t, err)
	assert.Equal(t, "оригинальный URL уже сопоставлен с другим коротким URL", err.Error())
}

func TestMemoryStorage_GetNotFound(t *testing.T) {
	store := memory.NewStorage()

	_, err := store.GetShort("https://nonexistent.com")
	assert.Error(t, err)

	_, err = store.GetOriginal("NonExistent")
	assert.Error(t, err)
}
