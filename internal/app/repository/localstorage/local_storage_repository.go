package localstorage

import (
	"context"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/jmoiron/sqlx"
)

type LocalStorageRepository struct {
	Rows map[string]LocalStorageRow
}

func (repo *LocalStorageRepository) GetUrlsByUserId(userId string) ([]entity.ShortURLEntity, error) {
	var models []entity.ShortURLEntity
	return models, nil
}

type LocalStorageRow struct {
	ShortURL      string
	URL           string
	CorrelationID string
	AddedUserId   string
}

func (repo *LocalStorageRepository) Batch(ents []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	var shortURL string
	var models []entity.ShortURLEntity
	for _, ent := range ents {
		shortURL = utils.GenerateShortURL(ent.OriginalURL)
		repo.Rows[shortURL] = LocalStorageRow{ShortURL: shortURL, URL: ent.OriginalURL, CorrelationID: ent.CorrelationID}
		models = append(models, entity.ShortURLEntity{ShortURL: shortURL, URL: ent.OriginalURL, CorrelationID: ent.CorrelationID})
	}

	return models, nil
}

func NewLocalStorageRepository() *LocalStorageRepository {
	return &LocalStorageRepository{
		Rows: make(map[string]LocalStorageRow),
	}
}

func (repo *LocalStorageRepository) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	if row, exists := repo.Rows[shortURL]; exists {
		return &(entity.ShortURLEntity{ShortURL: row.ShortURL, URL: row.URL, CorrelationID: row.CorrelationID}), nil
	}

	return nil, nil
}

func (repo *LocalStorageRepository) Save(shortURL, URL, addedUserId string) error {
	repo.Rows[shortURL] = LocalStorageRow{ShortURL: shortURL, URL: URL, AddedUserId: addedUserId}

	return nil
}

func (repo *LocalStorageRepository) Ping(ctx context.Context) error {

	return nil
}

func (repo *LocalStorageRepository) GetDB() *sqlx.DB {
	return nil
}
