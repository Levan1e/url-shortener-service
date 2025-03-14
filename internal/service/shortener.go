package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"
)

const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	shortURLLength = 10
	maxAttempts    = 5
)

type Storage interface {
	Save(originalURL, shortURL string) error

	GetShort(originalURL string) (string, error)

	GetOriginal(shortURL string) (string, error)
}

type SaveTask struct {
	originalURL string
	shortURL    string
	result      chan error
}

type AsyncSaver struct {
	storage Storage
	tasks   chan SaveTask
}

func NewAsyncSaver(storage Storage, workerCount int) *AsyncSaver {
	saver := &AsyncSaver{
		storage: storage,
		tasks:   make(chan SaveTask, 1000),
	}
	for i := 0; i < workerCount; i++ {
		go saver.worker()
	}
	return saver
}

func (a *AsyncSaver) worker() {
	for task := range a.tasks {
		err := a.storage.Save(task.originalURL, task.shortURL)
		task.result <- err
	}
}

func (a *AsyncSaver) SaveAsync(originalURL, shortURL string) error {
	result := make(chan error, 1)
	task := SaveTask{
		originalURL: originalURL,
		shortURL:    shortURL,
		result:      result,
	}
	a.tasks <- task
	select {
	case err := <-result:
		return err
	case <-time.After(5 * time.Second):
		return errors.New("timeout saving URL")
	}
}

type ShortenerService struct {
	storage    Storage
	asyncSaver *AsyncSaver
}

func NewShortenerService(storage Storage) *ShortenerService {
	return &ShortenerService{
		storage:    storage,
		asyncSaver: NewAsyncSaver(storage, 5),
	}
}

func (s *ShortenerService) Shorten(originalURL string) (string, error) {
	if short, err := s.storage.GetShort(originalURL); err == nil && short != "" {
		return short, nil
	}

	var shortURL string
	var err error

	for i := 0; i < maxAttempts; i++ {
		shortURL, err = generateRandomString(shortURLLength)
		if err != nil {
			return "", err
		}

		err = s.asyncSaver.SaveAsync(originalURL, shortURL)
		if err == nil {
			return shortURL, nil
		}
	}

	return "", errors.New("не удалось сгенерировать уникальный короткий URL")
}

func generateRandomString(n int) (string, error) {
	result := make([]byte, n)
	charsetLength := big.NewInt(int64(len(charset)))
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func (s *ShortenerService) Resolve(shortURL string) (string, error) {
	return s.storage.GetOriginal(shortURL)
}
