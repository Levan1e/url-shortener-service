package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Levan1e/url-shortener-service/internal/api"
	"github.com/Levan1e/url-shortener-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type stubStorage struct{}

func (s stubStorage) Save(originalURL, shortURL string) error {
	if originalURL == "error" {
		return errors.New("error saving")
	}
	return nil
}

func (s stubStorage) GetShort(originalURL string) (string, error) {
	if originalURL == "https://example.com" {
		return "AbCdEfGh12", nil
	}
	return "", errors.New("not found")
}

func (s stubStorage) GetOriginal(shortURL string) (string, error) {
	if shortURL == "AbCdEfGh12" {
		return "https://example.com", nil
	}
	return "", errors.New("not found")
}

func TestCreateShortURL_Success(t *testing.T) {
	storageStub := stubStorage{}
	svc := service.NewShortenerService(storageStub)
	handler := api.NewHandler(svc)

	reqBody, _ := json.Marshal(map[string]string{"url": "https://example.com"})
	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/shorten", handler.CreateShortURL)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "AbCdEfGh12", resp["short_url"])
}

func TestCreateShortURL_BadRequest(t *testing.T) {
	storageStub := stubStorage{}
	svc := service.NewShortenerService(storageStub)
	handler := api.NewHandler(svc)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/shorten", handler.CreateShortURL)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetOriginalURL_Success(t *testing.T) {
	storageStub := stubStorage{}
	svc := service.NewShortenerService(storageStub)
	handler := api.NewHandler(svc)

	req := httptest.NewRequest("GET", "/AbCdEfGh12", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/{shortURL}", handler.GetOriginalURL)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", resp["original_url"])
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	storageStub := stubStorage{}
	svc := service.NewShortenerService(storageStub)
	handler := api.NewHandler(svc)

	req := httptest.NewRequest("GET", "/notfound", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/{shortURL}", handler.GetOriginalURL)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
