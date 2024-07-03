package local_storage

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

func (repo *LocalStorageRepository) FindByShortUrl(shortUrl string) (*entity.ShortUrlEntity, error) {
	if url, exists := repo.Rows[shortUrl]; exists {
		return &(entity.ShortUrlEntity{ShortUrl: shortUrl, Url: url}), nil
	}

	return nil, fmt.Errorf("Row not found by short url: %s", shortUrl)
}

func (repo *LocalStorageRepository) Save(shortUrl string, url string) error {
	repo.Rows[shortUrl] = url

	return nil
}
