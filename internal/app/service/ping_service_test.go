package service

import (
	"context"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/jmoiron/sqlx"
	"testing"
)

type RepositoryMock struct{}

func (r RepositoryMock) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {
	return nil, nil
}
func (r RepositoryMock) Ping(ctx context.Context) error {
	return nil
}
func (r RepositoryMock) GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error) {
	return make([]entity.ShortURLEntity, 0), nil
}
func (r RepositoryMock) GetDB() *sqlx.DB {
	return nil
}

func TestNewPingService(t *testing.T) {
	repository := RepositoryMock{}
	tests := []struct {
		name string
	}{
		{
			name: "NewPingService",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewPingService(repository)
		})
	}
}
