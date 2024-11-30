package filestorage

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewFileStorageRepositoryWrite(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *FileStorageWriteRepository
		wantErr bool
	}{
		{
			name:    "success",
			path:    "/tmp/filestorage-test.txt",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFileStorageRepositoryWrite(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileStorageRepositoryWrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewFileStorageRepositoryRead(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *FileStorageWriteRepository
		wantErr bool
	}{
		{
			name:    "success",
			path:    "/tmp/filestorage-test.txt",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFileStorageRepositoryRead(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileStorageRepositoryRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFileStorageReadRepository_FindByShortURL(t *testing.T) {
	repositoryRead, err := NewFileStorageRepositoryRead("/tmp/filestorage-test.txt")
	if err != nil {
		t.Fatal(err)
	}
	repositoryWrite, err := NewFileStorageRepositoryWrite("/tmp/filestorage-test.txt")
	if err != nil {
		t.Fatal(err)
	}

	repositoryWrite.Save("bla", "bla", "bla")

	tests := []struct {
		name         string
		findShortURL string
		want         *entity.ShortURLEntity
		wantErr      bool
	}{
		{
			name:         "not find",
			findShortURL: "bla2",
			want:         nil,
			wantErr:      false,
		},

		{
			name:         "find",
			findShortURL: "bla",
			want: &entity.ShortURLEntity{
				ShortURL:      "bla",
				URL:           "bla",
				CorrelationID: "",
				AddedUserID:   "bla",
				IsDeleted:     false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repositoryRead.FindByShortURL(tt.findShortURL)
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

func TestFileStorageWriteRepository_Save(t *testing.T) {
	repositoryWrite, err := NewFileStorageRepositoryWrite("/tmp/filestorage-test.txt")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		shortURL    string
		URL         string
		addedUserID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				shortURL:    "bla",
				URL:         "bla",
				addedUserID: "bla",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repositoryWrite.Save(tt.args.shortURL, tt.args.URL, tt.args.addedUserID); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileStorageWriteRepository_Batch(t *testing.T) {
	repositoryWrite, err := NewFileStorageRepositoryWrite("/tmp/filestorage-test.txt")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		args    []batch.ButchRequest
		want    []entity.ShortURLEntity
		wantErr bool
	}{
		{
			name: "success",
			args: []batch.ButchRequest{
				{
					CorrelationID: "bla",
					OriginalURL:   "bla",
				},
				{
					CorrelationID: "foo",
					OriginalURL:   "foo",
				},
			},
			want: []entity.ShortURLEntity{
				{
					ShortURL:      "_6Zwb_IS",
					URL:           "bla",
					CorrelationID: "bla",
					AddedUserID:   "",
					IsDeleted:     false,
				},
				{
					ShortURL:      "C-7Hteo_",
					URL:           "foo",
					CorrelationID: "foo",
					AddedUserID:   "",
					IsDeleted:     false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repositoryWrite.Batch(tt.args)
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

func TestFileStorageReadRepository_Find(t *testing.T) {
	repositoryWrite, err := NewFileStorageRepositoryRead("/tmp/filestorage-test.txt")
	if err != nil {
		t.Fatal(err)
	}

	repositoryWrite.Find(nil, "bla")
}

func TestFileStorageReadRepository_GetUrlsByUserID(t *testing.T) {
	repositoryWrite, err := NewFileStorageRepositoryRead("/tmp/filestorage-test.txt")
	if err != nil {
		t.Fatal(err)
	}

	repositoryWrite.GetUrlsByUserID("bla")
}

func TestFileStorageReadRepository_GetDB(t *testing.T) {
	repositoryWrite, err := NewFileStorageRepositoryRead("/tmp/filestorage-test.txt")
	if err != nil {
		t.Fatal(err)
	}

	repositoryWrite.GetDB()
}

func TestFileStorageReadRepository_GetDB2(t *testing.T) {
	repositoryWrite, err := NewFileStorageRepositoryWrite("/tmp/filestorage-test.txt")
	if err != nil {
		t.Fatal(err)
	}

	repositoryWrite.GetDB()
}

func TestFileStorageReadRepository_Ping(t *testing.T) {
	repositoryWrite, err := NewFileStorageRepositoryRead("/tmp/filestorage-test.txt")
	if err != nil {
		t.Fatal(err)
	}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	repositoryWrite.Ping(ctx)
}
