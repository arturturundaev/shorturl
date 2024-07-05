package localstorage

import (
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
)

type LocalStorageRepository struct {
	Rows map[string]string
}

func NewLocalStorageRepository() *LocalStorageRepository {
	return &LocalStorageRepository{
		Rows: make(map[string]string),
	}
}

func (repo *LocalStorageRepository) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	if url, exists := repo.Rows[shortURL]; exists {
		return &(entity.ShortURLEntity{ShortURL: shortURL, URL: url}), nil
	}

	return nil, fmt.Errorf("row not found by short url: %s", shortURL)
}

func (repo *LocalStorageRepository) Save(shortURL string, URL string) error {
	repo.Rows[shortURL] = URL

	return nil
}
