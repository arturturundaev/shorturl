package main

import (
	"github.com/arturturundaev/shorturl/internal/app/repository/filestorage"
	"github.com/arturturundaev/shorturl/internal/app/repository/localstorage"
	"github.com/arturturundaev/shorturl/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"go.uber.org/zap"
	"net/http/httptest"
	"testing"
)

func Test_addLogger(t *testing.T) {
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	tests := []struct {
		name       string
		fullLogger bool
		want       *zap.Logger
		wantErr    bool
	}{
		{
			name:       "",
			fullLogger: false,
			want:       nil,
			wantErr:    false,
		},
		{
			name:       "",
			fullLogger: true,
			want:       nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := addLogger(engine, tt.fullLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("addLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_getRepository(t *testing.T) {
	tests := []struct {
		name         string
		serverConfig *config.Config
	}{
		{
			name:         "success connect to file storage",
			serverConfig: &config.Config{StorageType: config.StorageTypeFile, FileStorage: "/tmp/bla.txt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryReader, repositoryWriter := getRepository(tt.serverConfig, nil)

			_, ok := repositoryReader.(*filestorage.FileStorageReadRepository)
			_, ok2 := repositoryWriter.(*filestorage.FileStorageWriteRepository)

			if !ok {
				t.Errorf("getRepository() got = %v, want %v", repositoryWriter, nil)
			}
			if !ok2 {
				t.Errorf("getRepository() got1 = %v, want %v", repositoryReader, nil)
			}
		})
	}
}

func Test_getRepository2(t *testing.T) {
	tests := []struct {
		name         string
		serverConfig *config.Config
	}{
		{
			name:         "success connect to memory storage",
			serverConfig: &config.Config{StorageType: config.StorageTypeMemory},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryReader, repositoryWriter := getRepository(tt.serverConfig, nil)

			_, ok := repositoryReader.(*localstorage.LocalStorageRepository)
			_, ok2 := repositoryWriter.(*localstorage.LocalStorageRepository)

			if !ok {
				t.Errorf("getRepository() got = %v, want %v", repositoryWriter, nil)
			}
			if !ok2 {
				t.Errorf("getRepository() got1 = %v, want %v", repositoryReader, nil)
			}
		})
	}
}
func Test_initMigrations(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "creating a db mock")
	defer db.Close()
	//dbx := sqlx.NewDb(db, "sqlmock")
	mock.ExpectQuery("SELECT CURRENT_DATABASE()").WillReturnRows(mock.NewRows([]string{"b"}).AddRow("sqlmock"))
	mock.ExpectQuery("SELECT CURRENT_SCHEMA()").WillReturnRows(mock.NewRows([]string{"b"}).AddRow("public"))
	type args struct {
		migrationPath string
		DB            *sqlx.DB
	}
	tests := []struct {
		name string
		args args
	}{
		/*{
			name: "test",
			args: args{
				migrationPath: "/tmp/",
				DB:            dbx,
			},
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initMigrations(tt.args.migrationPath, tt.args.DB)
		})
	}
}
