package service

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/utils"
)

type ShortURLService struct {
	repository RepositoryInterface
}

func NewShortURLService(repository RepositoryInterface) *ShortURLService {
	return &ShortURLService{repository: repository}
}

func (service *ShortURLService) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	return service.repository.FindByShortURL(shortURL)
}

func (service *ShortURLService) Save(url string) (*entity.ShortURLEntity, error) {
	shortURL := utils.GenerateShortURL(url)

	model, _ := service.repository.FindByShortURL(shortURL)

	if model != nil {
		return model, nil
	}

	err := service.repository.Save(shortURL, url)

	if err != nil {
		return nil, err
	}

	return &entity.ShortURLEntity{ShortURL: shortURL, URL: url}, nil
}
