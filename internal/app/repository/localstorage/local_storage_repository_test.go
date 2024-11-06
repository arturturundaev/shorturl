package localstorage

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"reflect"
	"testing"
)

func TestLocalStorageRepository_Batch(t *testing.T) {
	tests := []struct {
		name    string
		args    []batch.ButchRequest
		want    []entity.ShortURLEntity
		wantErr bool
	}{
		{
			name: "Success count 2",
			args: []batch.ButchRequest{
				{
					CorrelationID: "1",
					OriginalURL:   "1",
				},
				{
					CorrelationID: "2",
					OriginalURL:   "2",
				},
			},
			want: []entity.ShortURLEntity{
				{
					ShortURL:      utils.GenerateShortURL("1"),
					URL:           "1",
					CorrelationID: "1",
					IsDeleted:     false,
				},
				{
					ShortURL:      utils.GenerateShortURL("2"),
					URL:           "2",
					CorrelationID: "2",
					IsDeleted:     false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewLocalStorageRepository()
			got, err := repo.Batch(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Batch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Batch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
