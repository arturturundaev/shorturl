package service

import "github.com/arturturundaev/shorturl/internal/app/entity"

type RepositoryReadInterface interface {
	FindByShortURL(shortURL string) (*entity.ShortURLEntity, error)
}

type RepositoryWriteInterface interface {
	Save(shortURL string, url string) error
}
