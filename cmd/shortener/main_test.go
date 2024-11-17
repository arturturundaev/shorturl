package main

import (
	"github.com/arturturundaev/shorturl/internal/app/repository/filestorage"
	"github.com/arturturundaev/shorturl/internal/app/repository/localstorage"
	"github.com/arturturundaev/shorturl/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http/httptest"
	"testing"
)

func Test_initRouter(t *testing.T) {
	tests := []struct {
		name  string
		want  *gin.Engine
		want1 *zap.Logger
		want2 *config.Config
	}{
		{
			name:  "success",
			want:  nil,
			want1: nil,
			want2: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _ = initRouter()
		})
	}
}

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
			serverConfig: &config.Config{StorageType: config.StorageTypeFile, FileStorage: config.FileStorageType{Path: "/tmp/bla.txt"}},
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
	type args struct {
		migrationPath string
		DB            *sqlx.DB
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initMigrations(tt.args.migrationPath, tt.args.DB)
		})
	}
}
