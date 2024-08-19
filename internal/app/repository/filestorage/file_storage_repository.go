package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/jmoiron/sqlx"
	"io"
	"os"
	"strings"
)

type FileStorageReadRepository struct {
	file *os.File
}

func (repo *FileStorageReadRepository) GetUrlsByUserId(userId string) ([]entity.ShortURLEntity, error) {
	var models []entity.ShortURLEntity

	return models, nil
}

func (repo *FileStorageReadRepository) GetDB() *sqlx.DB {
	return nil
}

func (repo *FileStorageWriteRepository) GetDB() *sqlx.DB {
	return nil
}

func (repo *FileStorageReadRepository) Ping(ctx context.Context) error {
	return nil
}

type FileStorageWriteRepository struct {
	file *os.File
}

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

func NewFileStorageRepositoryWrite(path string) (*FileStorageWriteRepository, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileStorageWriteRepository{file: file}, nil
}

func NewFileStorageRepositoryRead(path string) (*FileStorageReadRepository, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileStorageReadRepository{file: file}, nil
}

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

func (repo *FileStorageWriteRepository) Save(shortURL, URL, addedUserId string) error {
	_, err := repo.file.WriteString(fmt.Sprintf(`{"short_url":"%s","original_url":"%s", "added_user_id":"%s"}`+"\n", shortURL, URL, addedUserId))

	if err != nil {
		return err
	}

	return nil
}
