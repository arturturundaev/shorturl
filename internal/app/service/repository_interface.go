package service

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/jmoiron/sqlx"
)

type RepositoryReadInterface interface {
	FindByShortURL(shortURL string) (*entity.ShortURLEntity, error)
	Ping() error
	GetDB() *sqlx.DB
}

type RepositoryWriteInterface interface {
	Save(shortURL string, url string) error
	Batch(*[]entity.ShortURLEntity) error
	GetDB() *sqlx.DB
	BeginTransaction() error
	RollbackTransaction() error
	CommitTransaction() error
}
