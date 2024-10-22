package shorten

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewShortenRequest(t *testing.T) {

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	tests := []struct {
		name    string
		body    string
		want    *ShortenRequest
		wantErr error
	}{
		{
			name:    "empty url",
			body:    "",
			want:    &ShortenRequest{},
			wantErr: fmt.Errorf("EOF"),
		},
		{
			name:    "success",
			body:    `{"url":"https://google.com"}`,
			want:    &ShortenRequest{URL: "https://google.com"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resquest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			ctx.Request = resquest
			got, err := NewShortenRequest(ctx)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
