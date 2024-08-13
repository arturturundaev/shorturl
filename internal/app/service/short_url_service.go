package service

import (
	"errors"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/repository/postgres"
	"github.com/arturturundaev/shorturl/internal/app/utils"
)

type ShortURLService struct {
	repositoryRead  RepositoryReadInterface
	repositoryWrite RepositoryWriteInterface
}

var EntityExistsError = errors.New("Entity exists")

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
		return model, EntityExistsError
	}

	err := service.repositoryWrite.Save(shortURL, url)

	if err != nil {
		return nil, err
	}

	return &entity.ShortURLEntity{ShortURL: shortURL, URL: url}, nil
}

func (service *ShortURLService) Batch(request *[]batch.ButchRequest) (*[]entity.ShortURLEntity, error) {
	var models []entity.ShortURLEntity
	var allModels []entity.ShortURLEntity

	err := service.repositoryWrite.BeginTransaction()

	if err != nil {
		return nil, err
	}

	for i, item := range *request {
		models = append(models, entity.ShortURLEntity{URL: item.OriginalUrl, CorrelationId: item.CorrelationId, ShortURL: utils.GenerateShortURL(item.OriginalUrl)})

		if len(models) == postgres.ButchSize || len(*request) == i+1 {
			err = service.repositoryWrite.Batch(&models)
			if err != nil {
				err2 := service.repositoryWrite.RollbackTransaction()
				if err2 != nil {
					return nil, err2
				}
				return nil, err
			}
			err = service.repositoryWrite.CommitTransaction()
			if err != nil {
				return nil, err
			}
			allModels = append(allModels, models...)
			models = nil
		}
	}

	return &allModels, nil
}
