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
				AddressStart: "127.0.0.1:8086",
				BaseShort:    "127.0.0.1:8080",
				FileStorage:  "",
				DatabaseURL:  "",
				StorageType:  "Memory",
				FullLog:      true,
				HTTPS: struct {
					Enabled    bool `json:"enable_https"`
					SSLKeyPath string
					SSLPemPath string
				}(struct {
					Enabled    bool
					SSLKeyPath string
					SSLPemPath string
				}{Enabled: false, SSLKeyPath: "./auto_server.key", SSLPemPath: "./auto_server.pem"}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfig()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
