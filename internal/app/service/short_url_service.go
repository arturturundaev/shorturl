package service

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/utils"
)

type ShortURLService struct {
	repositoryRead  RepositoryReadInterface
	repositoryWrite RepositoryWriteInterface
}

func NewShortURLService(repositoryRead RepositoryReadInterface, repositoryWrite RepositoryWriteInterface) *ShortURLService {
	return &ShortURLService{repositoryRead: repositoryRead, repositoryWrite: repositoryWrite}
}

func (service *ShortURLService) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	return service.repositoryRead.FindByShortURL(shortURL)
}

func (service *ShortURLService) Save(url string) (*entity.ShortURLEntity, error) {
	shortURL := utils.GenerateShortURL(url)

	model, errRepository := service.repositoryRead.FindByShortURL(shortURL)

	if errRepository != nil {
		return nil, errRepository
	}

	if model != nil {
		return model, nil
	}

	err := service.repositoryWrite.Save(shortURL, url)

	if err != nil {
		return nil, err
	}

	return &entity.ShortURLEntity{ShortURL: shortURL, URL: url}, nil
}
