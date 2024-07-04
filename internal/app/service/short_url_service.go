package service

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/utils"
)

type ShortUrlService struct {
	repository RepositoryInterface
}

func NewShortUrlService(repository RepositoryInterface) *ShortUrlService {
	return &ShortUrlService{repository: repository}
}

func (service *ShortUrlService) FindByShortUrl(shortUrl string) (*entity.ShortUrlEntity, error) {
	return service.repository.FindByShortUrl(shortUrl)
}

func (service *ShortUrlService) Save(url string) (*entity.ShortUrlEntity, error) {
	shortUrl := utils.GenerateShortUrl(url)

	model, _ := service.repository.FindByShortUrl(shortUrl)

	if model != nil {
		return model, nil
	}

	err := service.repository.Save(shortUrl, url)

	if err != nil {
		return nil, err
	}

	return &entity.ShortUrlEntity{ShortUrl: shortUrl, Url: url}, nil
}
