package batch

import (
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type serviceURLButcherMock struct {
	mock.Mock
}

func (m *serviceURLButcherMock) Batch(request []ButchRequest) ([]entity.ShortURLEntity, error) {
	var data interface{}
	if request[0].CorrelationID == "service_batch_return_error" {
		data = "service_batch_return_error"
	} else if request[0].CorrelationID == "success" {
		data = "success"
	} else {
		data = request
	}
	args := m.Called(data)
	return args.Get(0).([]entity.ShortURLEntity), args.Error(1)
}

func TestButchHandler_Handle(t *testing.T) {
	type fields struct {
		serviceResponse struct {
			ents []entity.ShortURLEntity
			err  error
		}
		baseURL string
	}

	type expected struct {
		code     int
		response interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		request string
		expctd  expected
	}{
		{
			name:    "empty_request",
			request: "",
			fields: fields{
				serviceResponse: struct {
					ents []entity.ShortURLEntity
					err  error
				}{nil, fmt.Errorf("some error")},
				baseURL: "",
			},
			expctd: expected{
				code:     http.StatusBadRequest,
				response: gin.H{"error": fmt.Errorf("invalid request")},
			},
		},

		{
			name:    "service_batch_return_error",
			request: `[{"correlation_id":"service_batch_return_error","original_url":"sd"}]`,
			fields: fields{
				serviceResponse: struct {
					ents []entity.ShortURLEntity
					err  error
				}{nil, fmt.Errorf("some error")},
				baseURL: "",
			},
			expctd: expected{
				code:     http.StatusBadRequest,
				response: gin.H{"error": fmt.Errorf("invalid request")},
			},
		},

		{
			name:    "success",
			request: `[{"correlation_id":"success","original_url":"sd"}]`,
			fields: fields{
				serviceResponse: struct {
					ents []entity.ShortURLEntity
					err  error
				}{[]entity.ShortURLEntity{{
					ShortURL:      "ShortURL",
					URL:           "URL",
					CorrelationID: "CorrelationID",
					AddedUserID:   "AddedUserID",
					IsDeleted:     false,
				}}, nil},
				baseURL: "/bla/",
			},
			expctd: expected{
				code:     http.StatusCreated,
				response: gin.H{"error": fmt.Errorf("invalid request")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.request))
			mockService := new(serviceURLButcherMock)
			mockService.On("Batch", tt.name).Return(tt.fields.serviceResponse.ents, tt.fields.serviceResponse.err)
			h := &ButchHandler{
				service: mockService,
				baseURL: tt.fields.baseURL,
			}
			h.Handle(ctx)

			assert.Equal(t, tt.expctd.code, ctx.Writer.Status())
		})
	}
}
