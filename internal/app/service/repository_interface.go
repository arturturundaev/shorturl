package service

import (
	"context"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/jmoiron/sqlx"
)

type RepositoryReader interface {
	FindByShortURL(shortURL string) (*entity.ShortURLEntity, error)
	Ping(ctx context.Context) error
	GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error)
	GetDB() *sqlx.DB
}

type RepositoryWriter interface {
	Save(shortURL, url, addedUserID string) error
	Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error)
	GetDB() *sqlx.DB
	Delete(shortURLs []string, addedUserID string) error
}
