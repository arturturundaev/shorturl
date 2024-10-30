package ping

import (
	"context"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.service.On("Ping", tt.name).Return(tt.serviceReposnse)
			h := &PingHandler{
				service: tt.service,
			}

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			h.Handle(ctx)

			assert.Equal(t, tt.wuntCode, ctx.Writer.Status())
		})
	}
}
