package stats

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockService struct {
}

func (s MockService) GetUrlsAndUsersStat() (int32, int32) {
	return 1, 1
}
func TestUrlsAndUsersStatHandler_Handle(t *testing.T) {
	service := new(MockService)

	type response struct {
		Urls  int32 `json:"urls"`
		Users int32 `json:"users"`
	}

	tests := []struct {
		name string
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.JSON(http.StatusOK, response{
				Urls:  1,
				Users: 1,
			})
			h := NewUrlsAndUsersStatHandler(service)

			h.Handle(ctx)
		})
	}
}
