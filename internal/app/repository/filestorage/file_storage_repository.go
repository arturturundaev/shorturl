package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/jmoiron/sqlx"
	"io"
	"os"
	"strings"
)

type FileStorageReadRepository struct {
	file *os.File
}

func (repo *FileStorageReadRepository) GetDB() *sqlx.DB {
	return nil
}

func (repo *FileStorageReadRepository) Ping() error {
	return nil
}

type FileStorageWriteRepository struct {
	file *os.File
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

func (repo *FileStorageWriteRepository) Save(shortURL string, URL string) error {
	_, err := repo.file.WriteString(fmt.Sprintf(`{"short_url":"%s","original_url":"%s"}`+"\n", shortURL, URL))

	if err != nil {
		return err
	}

	return nil
}
