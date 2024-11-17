package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			"wthout any settings",
			&Config{
				AddressStart: AddressStartType{
					URL:  "localhost",
					Port: "8080",
				},
				BaseShort: BaseShortURLType{
					URL: "http://localhost:8080",
				},
				FileStorage: FileStorageType{
					Path: "/tmp/db.txt",
				},
				DatabaseURL: DatabaseURLType{
					URL: "postgres://postgres:postgres@localhost:5432/shorturl?sslmode=disable",
				},
				StorageType: "Memory",
				FullLog:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfig()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
			got.BaseShort.Set("s")
			got.FileStorage.Set("s")
			got.DatabaseURL.Set("s")
			got.AddressStart.Set("127.0.0.1:8080")
		})
	}
}
