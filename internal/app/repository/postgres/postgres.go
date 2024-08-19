package postgres

import (
	"context"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strings"
)

const TableName = "url"
const ButchSize = 100

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(databaseURL string) (*PostgresRepository, error) {
	database, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{DB: database}, nil
}

func (repo *PostgresRepository) Ping(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("прервали работу")
	default:
		return repo.DB.Ping()
	}
}

func (repo *PostgresRepository) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {

	ent := []entity.ShortURLEntity{}
	err := repo.DB.Select(&ent,
		fmt.Sprintf("select url_full as original_url, url_short from %s where url_short = $1", TableName),
		shortURL)

	if err != nil {
		return nil, err
	}

	if len(ent) == 0 {
		return nil, nil
	}

	return &ent[0], nil
}

func (repo *PostgresRepository) Save(shortURL, URL, addedUserID string) error {

	id := uuid.New().String()
	_, err := repo.DB.Exec(fmt.Sprintf(`INSERT into %s (id, url_full, url_short, added_user_id) values ($1, $2, $3, $4)`, TableName), id, URL, shortURL, addedUserID)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresRepository) GetDB() *sqlx.DB {
	return repo.DB
}

func (repo *PostgresRepository) Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	var models []entity.ShortURLEntity
	var allModels []entity.ShortURLEntity

	tx, err := repo.DB.Begin()

	if err != nil {
		return nil, err
	}

	k := 0
	var values []string
	var params []interface{}

	for i, item := range request {
		models = append(models, entity.ShortURLEntity{URL: item.OriginalURL, CorrelationID: item.CorrelationID, ShortURL: utils.GenerateShortURL(item.OriginalURL)})

		if len(models) == ButchSize || len(request) == i+1 {
			values = nil
			params = nil
			for _, enty := range models {
				values = append(values, "("+
					"$"+fmt.Sprintf("%d,", k+1)+
					"$"+fmt.Sprintf("%d,", k+2)+
					"$"+fmt.Sprintf("%d,", k+3)+
					"$"+fmt.Sprintf("%d", k+4)+")")
				params = append(params, uuid.New().String(), enty.URL, enty.ShortURL, enty.CorrelationID)
				k += 4
			}

			valuesStr := strings.Join(values, ",")
			_, err = repo.DB.Exec(fmt.Sprintf(`INSERT into %s values %s`, TableName, valuesStr), params...)
			if err != nil {
				err2 := tx.Rollback()
				if err2 != nil {
					return nil, err2
				}
				return nil, err
			}
			allModels = append(allModels, models...)
			models = nil
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return allModels, nil
}

func (repo *PostgresRepository) GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error) {
	ent := []entity.ShortURLEntity{}

	err := repo.DB.Select(&ent,
		fmt.Sprintf("select url_full as original_url, url_short from %s where added_user_id = $1", TableName),
		userID)

	if err != nil {
		return nil, err
	}

	return ent, nil
}
