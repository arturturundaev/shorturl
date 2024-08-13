package localstorage

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/jmoiron/sqlx"
)

type LocalStorageRepository struct {
	Rows map[string]LocalStorageRow
}

type LocalStorageRow struct {
	ShortURL      string
	URL           string
	CorrelationId string
}

func (repo *LocalStorageRepository) Batch(ents *[]entity.ShortURLEntity) error {
	for _, ent := range *ents {
		repo.Rows[ent.ShortURL] = LocalStorageRow{ShortURL: ent.ShortURL, URL: ent.URL, CorrelationId: ent.CorrelationId}
	}

	return nil
}

func NewLocalStorageRepository() *LocalStorageRepository {
	return &LocalStorageRepository{
		Rows: make(map[string]LocalStorageRow),
	}
}

func (repo *LocalStorageRepository) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	if row, exists := repo.Rows[shortURL]; exists {
		return &(entity.ShortURLEntity{ShortURL: row.ShortURL, URL: row.URL, CorrelationId: row.CorrelationId}), nil
	}

	return nil, nil
}

func (repo *LocalStorageRepository) Save(shortURL string, URL string) error {
	repo.Rows[shortURL] = LocalStorageRow{ShortURL: shortURL, URL: URL}

	return nil
}

func (repo *LocalStorageRepository) Ping() error {

	return nil
}

func (repo *LocalStorageRepository) GetDB() *sqlx.DB {
	return nil
}

func (repo *LocalStorageRepository) BeginTransaction() error {
	return nil
}

func (repo *LocalStorageRepository) RollbackTransaction() error {
	return nil
}

func (repo *LocalStorageRepository) CommitTransaction() error {
	return nil
}
