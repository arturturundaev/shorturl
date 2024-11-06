package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// TableName имя таблицы
const TableName = "url"

// Количество записей на одну вставвку
const ButchSize = 100

// PostgresRepository сервис
type PostgresRepository struct {
	DB *sqlx.DB
}

// NewPostgresRepository конструктор
func NewPostgresRepository(databaseURL string) (*PostgresRepository, error) {
	database, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{DB: database}, nil
}

// Ping  пинг
func (repo *PostgresRepository) Ping(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("прервали работу")
	default:
		return repo.DB.Ping()
	}
}

// FindByShortURL поиск по короткой ссылке
func (repo *PostgresRepository) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {

	ent := []entity.ShortURLEntity{}
	err := repo.DB.Select(&ent,
		fmt.Sprintf("select url_full as original_url, url_short, is_deleted from %s where url_short = $1", TableName),
		shortURL)

	if err != nil {
		return nil, err
	}

	if len(ent) == 0 {
		return nil, nil
	}

	return &ent[0], nil
}

// Save сохранение
func (repo *PostgresRepository) Save(shortURL, URL, addedUserID string) error {

	id := uuid.New().String()
	_, err := repo.DB.Exec(fmt.Sprintf(`INSERT into %s (id, url_full, url_short, added_user_id) values ($1, $2, $3, $4)`, TableName), id, URL, shortURL, addedUserID)

	if err != nil {
		return err
	}

	return nil
}

// GetDB получение коннекта к репозиторию
func (repo *PostgresRepository) GetDB() *sqlx.DB {
	return repo.DB
}

// Batch Массовое сохранение
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

// GetUrlsByUserID получение ссылок по пользователю
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

// Delete удаление
func (repo *PostgresRepository) Delete(shortURLs []string, addedUserID string) error {
	var inArray []string
	var params []interface{}

	params = append(params, addedUserID)
	i := 2
	for _, shortURL := range shortURLs {
		inArray = append(inArray, "$"+fmt.Sprintf("%d", i))
		params = append(params, shortURL)
		i++
	}
	_, err := repo.DB.Exec(fmt.Sprintf(`update %s SET is_deleted = true WHERE added_user_id = $1 AND url_short IN (%s)`, TableName, strings.Join(inArray, ",")), params...)

	if err != nil {
		return err
	}

	return nil
}
