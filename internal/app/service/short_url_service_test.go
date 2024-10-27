package service

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"testing"
)

type RepositoryWriteMock struct {
	mock.Mock
}

func (r RepositoryWriteMock) Save(shortURL, url, addedUserID string) error {
	return nil
}

func (r RepositoryWriteMock) Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	return make([]entity.ShortURLEntity, 0), nil
}

func (r RepositoryWriteMock) GetDB() *sqlx.DB {
	return nil
}

func (r RepositoryWriteMock) Delete(shortURLs []string, addedUserID string) error {
	return nil
}

func TestShortURLService_Delete(t *testing.T) {

	repositoryWrite := new(RepositoryWriteMock)
	type args struct {
		URLList     []string
		addedUserID string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				URLList:     []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
				addedUserID: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &ShortURLService{
				repositoryRead:  nil,
				repositoryWrite: repositoryWrite,
				logger:          nil,
			}
			service.Delete(tt.args.URLList, tt.args.addedUserID)
		})
	}
}
