package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/jmoiron/sqlx"
)

// FileStorageReadRepository сервис
type FileStorageReadRepository struct {
	file *os.File
}

// Find поиск
func (repo *FileStorageReadRepository) Find(shortURLs []string, addedUserID string) ([]entity.ShortURLEntity, error) {
	return make([]entity.ShortURLEntity, 0), nil
}

// GetUrlsByUserID получение ссылок по пользователю
func (repo *FileStorageReadRepository) GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error) {
	var models []entity.ShortURLEntity

	return models, nil
}

// GetDB получение коннекта к репозиторию на чтение
func (repo *FileStorageReadRepository) GetDB() *sqlx.DB {
	return nil
}

// GetDB получение коннекта к репозиторию на запись
func (repo *FileStorageWriteRepository) GetDB() *sqlx.DB {
	return nil
}

// Ping  пинг
func (repo *FileStorageReadRepository) Ping(ctx context.Context) error {
	return nil
}

// FileStorageWriteRepository сервис
type FileStorageWriteRepository struct {
	file *os.File
}

// Batch Массовое сохранение
func (repo *FileStorageWriteRepository) Batch(ents []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	var models []entity.ShortURLEntity
	var shortURL string
	for _, ent := range ents {
		shortURL = utils.GenerateShortURL(ent.OriginalURL)
		_, err := repo.file.WriteString(fmt.Sprintf(`{"short_url":"%s","original_url":"%s", "correlation_id":"%s"}`+"\n", shortURL, ent.OriginalURL, ent.CorrelationID))
		if err != nil {
			return nil, err
		}
		models = append(models, entity.ShortURLEntity{ShortURL: shortURL, URL: ent.OriginalURL, CorrelationID: ent.CorrelationID})
	}

	return models, nil
}

// NewFileStorageRepositoryWrite контсруктор на запись
func NewFileStorageRepositoryWrite(path string) (*FileStorageWriteRepository, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileStorageWriteRepository{file: file}, nil
}

// NewFileStorageRepositoryRead контсруктор на чтение
func NewFileStorageRepositoryRead(path string) (*FileStorageReadRepository, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileStorageReadRepository{file: file}, nil
}

// FindByShortURL поиск по короткой ссылке
func (repo *FileStorageReadRepository) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	var dto entity.ShortURLEntity

	_, err := repo.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(repo.file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), shortURL) {
			err := json.Unmarshal(scanner.Bytes(), &dto)
			if err != nil {
				return nil, err
			}

			return &dto, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

// Save сохранение
func (repo *FileStorageWriteRepository) Save(shortURL, URL, addedUserID string) error {
	_, err := repo.file.WriteString(fmt.Sprintf(`{"short_url":"%s","original_url":"%s", "added_user_id":"%s"}`+"\n", shortURL, URL, addedUserID))

	if err != nil {
		return err
	}

	return nil
}

// Delete удаление
func (repo *FileStorageWriteRepository) Delete(shortURLs []string, addedUserID string) error {
	return nil
}
