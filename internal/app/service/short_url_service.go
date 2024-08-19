package service

import (
	"errors"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/gin-gonic/gin"
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

func (service *ShortURLService) Save(ctx *gin.Context, url string) (*entity.ShortURLEntity, error) {
	addedUserId, exists := ctx.Get(middleware.USER_ID_PROPERTY)

	if !exists {
		return nil, fmt.Errorf("user id is required")
	}

	shortURL := utils.GenerateShortURL(url)

	model, errRepository := service.repositoryRead.FindByShortURL(shortURL)

	if errRepository != nil {
		return nil, errRepository
	}

	if model != nil {
		return model, ErrEntityExists
	}

	err := service.repositoryWrite.Save(shortURL, url, addedUserId.(string))

	if err != nil {
		return nil, err
	}

	return &entity.ShortURLEntity{ShortURL: shortURL, URL: url}, nil
}

func (service *ShortURLService) Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	return service.repositoryWrite.Batch(request)
}

func (service *ShortURLService) GetUrlsByUserId(userId string) ([]entity.ShortURLEntity, error) {
	return service.repositoryRead.GetUrlsByUserId(userId)
}
