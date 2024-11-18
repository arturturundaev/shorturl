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

func TestLocalStorageRepository_FindByShortURL(t *testing.T) {
	repo := NewLocalStorageRepository()
	err := repo.Save("bla", "bla", "bla")
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}
	tests := []struct {
		name     string
		shortURL string
		want     *entity.ShortURLEntity
		wantErr  bool
	}{
		{
			name:     "not exists",
			shortURL: "bla2",
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "exists",
			shortURL: "bla",
			want: &entity.ShortURLEntity{
				ShortURL:      "bla",
				URL:           "bla",
				CorrelationID: "",
				AddedUserID:   "",
				IsDeleted:     false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByShortURL(tt.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByShortURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalStorageRepository_Save(t *testing.T) {
	repo := NewLocalStorageRepository()
	err := repo.Save("bla", "bla", "bla")
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}
}

func TestLocalStorageRepository_Ping(t *testing.T) {
	repo := NewLocalStorageRepository()
	resultPing := repo.Ping(nil)

	if resultPing != nil {
		t.Errorf("Ping() error")
	}
}

func TestLocalStorageRepository_GetDB(t *testing.T) {
	repo := NewLocalStorageRepository()
	db := repo.GetDB()

	if db != nil {
		t.Errorf("GetDB() error")
	}
}

func TestLocalStorageRepository_Delete(t *testing.T) {
	repo := NewLocalStorageRepository()
	result := repo.Delete([]string{"bla"}, "bla")

	if result != nil {
		t.Errorf("Delete() error")
	}
}

func TestLocalStorageRepository_GetUrlsByUserID(t *testing.T) {
	repo := NewLocalStorageRepository()
	got, err := repo.GetUrlsByUserID("bla")

	if err != nil {
		t.Errorf("Delete() error")
	}

	wnt := make([]entity.ShortURLEntity, 0)
	if reflect.DeepEqual(got, wnt) {
		t.Errorf("GetUrlsByUserID() error")
	}
}

func TestLocalStorageRepository_SaveToFile(t *testing.T) {
	repo := NewLocalStorageRepository()
	repo.Save("bla", "bla", "bla")
	repo.SaveToFile("/tmp/test.txt")
}
