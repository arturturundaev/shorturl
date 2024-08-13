package postgres

import (
	"database/sql"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strings"
)

const TableName = "url"
const ButchSize = 100

type PostgresRepository struct {
	DB *sqlx.DB
	tx *sql.Tx
}

func NewPostgresRepository(databaseURL string) (*PostgresRepository, error) {
	database, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{DB: database}, nil
}

func (repo *PostgresRepository) Ping() error {

	return repo.DB.Ping()
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

func (repo *PostgresRepository) Save(shortURL string, URL string) error {

	id := uuid.New().String()
	_, err := repo.DB.Exec(fmt.Sprintf(`INSERT into %s values ($1, $2, $3)`, TableName), id, URL, shortURL)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresRepository) GetDB() *sqlx.DB {
	return repo.DB
}

func (repo *PostgresRepository) Batch(entities *[]entity.ShortURLEntity) error {
	if len(*entities) == 0 {
		return nil
	}
	i := 0
	var values []string
	var params []interface{}
	for _, enty := range *entities {
		//$" + fmt.Sprintf("%d", paramIndex)
		values = append(values, "("+
			"$"+fmt.Sprintf("%d,", i+1)+
			"$"+fmt.Sprintf("%d,", i+2)+
			"$"+fmt.Sprintf("%d,", i+3)+
			"$"+fmt.Sprintf("%d", i+4)+")")
		params = append(params, uuid.New().String(), enty.URL, enty.ShortURL, enty.CorrelationID)
		i += 4
	}

	valuesStr := strings.Join(values, ",")
	_, err := repo.DB.Exec(fmt.Sprintf(`INSERT into %s values %s`, TableName, valuesStr), params...)

	return err
}

func (repo *PostgresRepository) BeginTransaction() error {
	tx, err := repo.DB.Begin()

	repo.tx = tx

	return err
}
func (repo *PostgresRepository) RollbackTransaction() error {
	return repo.tx.Rollback()
}
func (repo *PostgresRepository) CommitTransaction() error {
	return repo.tx.Commit()
}
