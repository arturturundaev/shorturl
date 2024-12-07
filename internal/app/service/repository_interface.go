package service

import (
	"context"

	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/jmoiron/sqlx"
)

// RepositoryReader интерфейс на запись
type RepositoryReader interface {
	FindByShortURL(shortURL string) (*entity.ShortURLEntity, error)
	Ping(ctx context.Context) error
	GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error)
	GetDB() *sqlx.DB
	GetUrlsCount() int32
	GetUsersCount() int32
}

// Интрефейс на чтение
type RepositoryWriter interface {
	Save(shortURL, url, addedUserID string) error
	Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error)
	GetDB() *sqlx.DB
	Delete(shortURLs []string, addedUserID string) error
	SaveToFile(fileName string) error
}
