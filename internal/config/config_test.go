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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
