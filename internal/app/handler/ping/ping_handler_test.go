package ping

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type pigngerMock struct {
	mock.Mock
}

func (p *pigngerMock) Ping(ctx context.Context) error {
	result := ctx.Value("result")

	if result == nil {
		return nil
	}

	if result == "error" {
		return fmt.Errorf("error")
	}
	return nil
}

func TestPingHandler_Handle(t *testing.T) {
	srv := new(pigngerMock)

	tests := []struct {
		name            string
		service         *pigngerMock
		serviceReposnse error
		wuntCode        int
	}{
		{
			name:            "success",
			service:         srv,
			serviceReposnse: nil,
			wuntCode:        http.StatusOK,
		},
		{
			name:            "error",
			service:         srv,
			serviceReposnse: fmt.Errorf("error"),
			wuntCode:        http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Set("result", tt.name)
			tt.service.On("Ping", ctx).Return(tt.serviceReposnse)
			h := &PingHandler{
				service: tt.service,
			}

			h.Handle(ctx)

			assert.Equal(t, tt.wuntCode, ctx.Writer.Status())
		})
	}
}
