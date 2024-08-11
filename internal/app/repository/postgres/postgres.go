package postgres

import (
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const TableName = "url"

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
