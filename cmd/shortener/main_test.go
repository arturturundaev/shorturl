package main

import (
	"github.com/arturturundaev/shorturl/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
