package service

import "github.com/arturturundaev/shorturl/internal/app/entity"

type RepositoryInterface interface {
	FindByShortUrl(shortUrl string) (*entity.ShortUrlEntity, error)
	Save(shortUrl string, url string) error
}
