package delete

import (
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type shortURLServiceMock struct {
	mock.Mock
}

func (service *shortURLServiceMock) Delete(URLList []string, addedUserID string) {

}

func TestDeleteHandler_Handle(t *testing.T) {

	service := new(shortURLServiceMock)
	handler := NewDeleteHandler(service)

	type testCase struct {
		name                   string
		requestBody            string
		userID                 string
		responseStatus         int
		serviceDeleteExecTimes int
	}

	tests := []testCase{
		{
			name:                   "Unmarshal fail",
			requestBody:            "none",
			userID:                 "none",
			responseStatus:         http.StatusBadRequest,
			serviceDeleteExecTimes: 0,
		},
		{
			name:                   "not auth",
			requestBody:            `["bla"]`,
			userID:                 "none",
			responseStatus:         http.StatusUnauthorized,
			serviceDeleteExecTimes: 0,
		},
		{
			name:                   "success",
			requestBody:            `["bla"]`,
			userID:                 "bla",
			responseStatus:         http.StatusAccepted,
			serviceDeleteExecTimes: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))

			if tt.userID != "none" {
				ctx.Set(middleware.UserIDProperty, tt.userID)
			}
			handler.Handle(ctx)
			assert.Equal(t, tt.responseStatus, ctx.Writer.Status())
			service.AssertNumberOfCalls(t, "Delete", tt.serviceDeleteExecTimes)
		})
	}
}
