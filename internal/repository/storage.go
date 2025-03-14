package repository

type Storage interface {
	Save(originalURL, shortURL string) error

	GetShort(originalURL string) (string, error)

	GetOriginal(shortURL string) (string, error)
}
