package config

import (
	"net"
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
				AddressStart: "localhost:8080",
				BaseShort:    "http://localhost:8080",
				FileStorage:  "FILE_STORAGE_PATH",
				DatabaseURL:  "DATABASE_DSN",
				StorageType:  "DB",
				FullLog:      true,
				HTTPS: struct {
					Enabled    bool `json:"enable_https"`
					SSLKeyPath string
					SSLPemPath string
				}(struct {
					Enabled    bool
					SSLKeyPath string
					SSLPemPath string
				}{Enabled: true, SSLKeyPath: "./auto_server.key", SSLPemPath: "./auto_server.pem"}),
				TrustedSubnet: "192.0.2.1/24",
				TrustedSubnetFinal: []*net.IPNet{{
					IP:   net.IPv4(192, 0, 2, 0).Mask(net.IPv4Mask(255, 255, 255, 0)),
					Mask: net.IPv4Mask(255, 255, 255, 0),
				},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("FILE_STORAGE_PATH", "FILE_STORAGE_PATH")
			t.Setenv("DATABASE_DSN", "DATABASE_DSN")
			t.Setenv("ENABLE_HTTPS", "true")
			t.Setenv("TRUSTED_SUBNET", "192.0.2.1/24")

			got := NewConfig()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
