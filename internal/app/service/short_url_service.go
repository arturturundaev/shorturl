package service

import (
	"errors"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/utils"
)

type ShortURLService struct {
	repositoryRead  RepositoryReader
	repositoryWrite RepositoryWriter
}

var ErrEntityExists = errors.New("entity exists")

func NewShortURLService(repositoryRead RepositoryReader, repositoryWrite RepositoryWriter) *ShortURLService {
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
		return model, ErrEntityExists
	}

	err := service.repositoryWrite.Save(shortURL, url)

	if err != nil {
		return nil, err
	}

	return &entity.ShortURLEntity{ShortURL: shortURL, URL: url}, nil
}

func (service *ShortURLService) Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	return service.repositoryWrite.Batch(request)
}
