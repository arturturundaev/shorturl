package service

import "github.com/arturturundaev/shorturl/internal/app/entity"

type RepositoryInterface interface {
	FindByShortURL(shortURL string) (*entity.ShortURLEntity, error)
	Save(shortURL string, url string) error
}
