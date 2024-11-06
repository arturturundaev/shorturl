package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ShortURLService сервис
type ShortURLService struct {
	repositoryRead  RepositoryReader
	repositoryWrite RepositoryWriter
	logger          *zap.Logger
}

// ErrEntityExists ошибка нет записи
var ErrEntityExists = errors.New("entity exists")

// NewShortURLService конструктор
func NewShortURLService(repositoryRead RepositoryReader, repositoryWrite RepositoryWriter) *ShortURLService {
	return &ShortURLService{repositoryRead: repositoryRead, repositoryWrite: repositoryWrite}
}

// FindByShortURL поиск по короткой ссылке
func (service *ShortURLService) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	return service.repositoryRead.FindByShortURL(shortURL)
}

// Save сохранение
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

// Batch Массовое сохранение
func (service *ShortURLService) Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	return service.repositoryWrite.Batch(request)
}

// GetUrlsByUserID получение ссылок по пользователю
func (service *ShortURLService) GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error) {
	return service.repositoryRead.GetUrlsByUserID(userID)
}

// Delete удаление
func (service *ShortURLService) Delete(URLList []string, addedUserID string) {
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
	ch1 := service.prepareGoodURL(inCh, addedUserID)
	ch2 := service.prepareGoodURL(inCh, addedUserID)
	for n := range service.fanIn(ch1, ch2) {
		deletedURLs = append(deletedURLs, n...)
	}

	service.notification(deletedURLs)
}

// sendToPrepare отправляем в поток на обработку
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
func (service *ShortURLService) prepareGoodURL(inCh chan []string, addedUserID string) chan []string {
	outCh := make(chan []string)

	go func() {
		defer close(outCh)
		for shortURLs := range inCh {
			err := service.repositoryWrite.Delete(shortURLs, addedUserID)
			if err != nil {
				service.logger.Error("ошибка получения записей", zap.String("urls", strings.Join(shortURLs, ",")), zap.String("addedUserID", addedUserID), zap.Error(err))
				return
			}

			outCh <- shortURLs
		}
	}()

	return outCh
}

// fanIn обработка потока
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

// notification уведомление о обработанных письмах
func (service *ShortURLService) notification(URLs []string) {
}
