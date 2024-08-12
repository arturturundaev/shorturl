package localstorage

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/jmoiron/sqlx"
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

	return nil, nil
}

func (repo *LocalStorageRepository) Save(shortURL string, URL string) error {
	repo.Rows[shortURL] = URL

	return nil
}

func (repo *LocalStorageRepository) Ping() error {

	return nil
}

func (repo *LocalStorageRepository) GetDB() *sqlx.DB {
	return nil
}
