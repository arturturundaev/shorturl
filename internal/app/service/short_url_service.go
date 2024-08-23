package service

import (
	"errors"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"sync"
)

type ShortURLService struct {
	repositoryRead  RepositoryReader
	repositoryWrite RepositoryWriter
	logger          *zap.Logger
}

var ErrEntityExists = errors.New("entity exists")

func NewShortURLService(repositoryRead RepositoryReader, repositoryWrite RepositoryWriter) *ShortURLService {
	return &ShortURLService{repositoryRead: repositoryRead, repositoryWrite: repositoryWrite}
}

func (service *ShortURLService) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	return service.repositoryRead.FindByShortURL(shortURL)
}

func (service *ShortURLService) Save(ctx *gin.Context, url string) (*entity.ShortURLEntity, error) {
	addedUserID, exists := ctx.Get(middleware.UserIDProperty)

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

	err := service.repositoryWrite.Save(shortURL, url, addedUserID.(string))

	if err != nil {
		return nil, err
	}

	return &entity.ShortURLEntity{ShortURL: shortURL, URL: url}, nil
}

func (service *ShortURLService) Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	return service.repositoryWrite.Batch(request)
}

func (service *ShortURLService) GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error) {
	return service.repositoryRead.GetUrlsByUserID(userID)
}

func (service *ShortURLService) Delete(URLList []string, addedUserId string) {
	service.logger.Info("Пытаемся удалить url: " + strings.Join(URLList, ",") + " addedUserId: " + addedUserId)
	var chunk []string
	var chunks [][]string
	i := 0
	for _, URL := range URLList {
		chunk = append(chunk, URL)

		if i < 10 {
			i++
		} else {
			chunks = append(chunks, chunk)
			chunk = nil
			i = 0
		}

	}

	var deletedURLs []string
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
		chunk = nil
	}

	inCh := service.sendToPrepare(chunks)
	ch1 := service.prepareGoodURL(inCh, addedUserId)
	ch2 := service.prepareGoodURL(inCh, addedUserId)
	for n := range service.fanIn(ch1, ch2) {
		deletedURLs = append(deletedURLs, n...)
	}

	service.notification(deletedURLs)
}

func (service *ShortURLService) sendToPrepare(chunks [][]string) chan []string {
	outCh := make(chan []string)
	go func() {
		defer close(outCh)
		for _, chunk := range chunks {
			outCh <- chunk
		}
	}()

	return outCh
}

// Возвращаем только те URL, которые действо можно удалить
func (service *ShortURLService) prepareGoodURL(inCh chan []string, addedUserId string) chan []string {
	outCh := make(chan []string)

	go func() {
		defer close(outCh)
		for shortURLs := range inCh {
			err := service.repositoryWrite.Delete(shortURLs, addedUserId)
			if err != nil {
				service.logger.Error("ошибка получения записей", zap.String("urls", strings.Join(shortURLs, ",")), zap.String("addedUserId", addedUserId), zap.Error(err))
				return
			}

			outCh <- shortURLs
		}
	}()

	return outCh
}

func (service *ShortURLService) fanIn(chs ...chan []string) chan []string {
	var wg sync.WaitGroup
	outCh := make(chan []string)

	// определяем функцию output для каждого канала в chs
	// функция output копирует значения из канала с в канал outCh, пока с не будет закрыт
	output := func(c chan []string) {
		for n := range c {
			outCh <- n
		}
		wg.Done()
	}

	// добавляем в группу столько горутин, сколько каналов пришло в fanIn
	wg.Add(len(chs))
	// перебираем все каналы, которые пришли и отправляем каждый в отдельную горутину
	for _, c := range chs {
		go output(c)
	}

	// запускаем горутину для закрытия outCh после того, как все горутины отработают
	go func() {
		wg.Wait()
		close(outCh)
	}()

	// возвращаем общий канал
	return outCh
}

func (service *ShortURLService) notification(URLs []string) {
}
