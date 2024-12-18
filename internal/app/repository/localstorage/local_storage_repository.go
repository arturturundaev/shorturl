package localstorage

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/jmoiron/sqlx"
	"os"
)

// LocalStorageRepository сервис
type LocalStorageRepository struct {
	Rows map[string]LocalStorageRow
}

// GetUrlsByUserID получение ссылок по пользователю
func (repo *LocalStorageRepository) GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error) {
	var models []entity.ShortURLEntity
	return models, nil
}

// LocalStorageRow сервис
type LocalStorageRow struct {
	ShortURL      string
	URL           string
	CorrelationID string
	AddedUserID   string
}

// Batch Массовое сохранение
func (repo *LocalStorageRepository) Batch(ents []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	var shortURL string
	var models []entity.ShortURLEntity
	for _, ent := range ents {
		shortURL = utils.GenerateShortURL(ent.OriginalURL)
		repo.Rows[shortURL] = LocalStorageRow{ShortURL: shortURL, URL: ent.OriginalURL, CorrelationID: ent.CorrelationID}
		models = append(models, entity.ShortURLEntity{ShortURL: shortURL, URL: ent.OriginalURL, CorrelationID: ent.CorrelationID})
	}

	return models, nil
}

// NewLocalStorageRepository конструктор
func NewLocalStorageRepository() *LocalStorageRepository {
	return &LocalStorageRepository{
		Rows: make(map[string]LocalStorageRow),
	}
}

// FindByShortURL поиск по короткой ссылке
func (repo *LocalStorageRepository) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	if row, exists := repo.Rows[shortURL]; exists {
		return &(entity.ShortURLEntity{ShortURL: row.ShortURL, URL: row.URL, CorrelationID: row.CorrelationID}), nil
	}

	return nil, nil
}

// Save сохранение
func (repo *LocalStorageRepository) Save(shortURL, URL, addedUserID string) error {
	repo.Rows[shortURL] = LocalStorageRow{ShortURL: shortURL, URL: URL, AddedUserID: addedUserID}

	return nil
}

// Ping  пинг
func (repo *LocalStorageRepository) Ping(ctx context.Context) error {

	return nil
}

// GetDB получение коннекта к репозиторию
func (repo *LocalStorageRepository) GetDB() *sqlx.DB {
	return nil
}

// Delete удаление
func (repo *LocalStorageRepository) Delete(shortURLs []string, addedUserID string) error {
	return nil
}

// SaveToFile сохраняем текущее состояние в файл при падении
func (repo *LocalStorageRepository) SaveToFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	for _, v := range repo.Rows {
		err = encoder.Encode(v)
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
