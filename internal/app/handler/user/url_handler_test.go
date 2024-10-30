package user

import (
	"errors"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type URLServiceFinderTest struct {
	mock.Mock
}

func (service URLServiceFinderTest) GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error) {
	args := service.Called(userID)

	return args.Get(0).([]entity.ShortURLEntity), args.Error(1)
}

func TestURLFindByUserHandler_Handle(t *testing.T) {

	service := new(URLServiceFinderTest)

	type serviceGetUrlsByUserID struct {
		entities []entity.ShortURLEntity
		err      error
	}

	type testCase struct {
		name            string
		userID          string
		serviceResponse serviceGetUrlsByUserID
		responseStatus  int
	}

	tests := []testCase{
		{
			name:   "fail find in repository",
			userID: "fail",
			serviceResponse: serviceGetUrlsByUserID{
				entities: nil,
				err:      errors.New("repository error"),
			},
			responseStatus: http.StatusBadRequest,
		},
		{
			name:   "no auth",
			userID: "no_auth",
			serviceResponse: serviceGetUrlsByUserID{
				entities: make([]entity.ShortURLEntity, 0),
				err:      nil,
			},
			responseStatus: http.StatusUnauthorized,
		},
		{
			name:   "success",
			userID: "success",
			serviceResponse: serviceGetUrlsByUserID{
				entities: []entity.ShortURLEntity{
					{
						ShortURL:      "success",
						URL:           "success",
						CorrelationID: "success",
						AddedUserID:   "success",
						IsDeleted:     false,
					},
				},
				err: nil,
			},
			responseStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Set(middleware.UserIDProperty, tt.userID)

			service.On("GetUrlsByUserID", tt.userID).Return(tt.serviceResponse.entities, tt.serviceResponse.err)
			handler := &URLFindByUserHandler{
				service: service,
				baseURL: "",
			}
			handler.Handle(ctx)
			assert.Equal(t, tt.responseStatus, ctx.Writer.Status())
		})
	}
}
