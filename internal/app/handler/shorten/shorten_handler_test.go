package shorten

import (
	"errors"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	service2 "github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ServiceMock struct{}

func (service *ServiceMock) Save(ctx *gin.Context, url string) (*entity.ShortURLEntity, error) {
	if url == "https://practicum.yandex.ru" {
		return &entity.ShortURLEntity{ShortURL: "7CwAhsKq", URL: "https://practicum.yandex.ru"}, nil
	}

	if url == "error" {
		return nil, errors.New("error")
	}
	if url == "exists" {
		return &entity.ShortURLEntity{ShortURL: "7CwAhsKq", URL: "https://practicum.yandex.ru"}, service2.ErrEntityExists
	}

	return nil, nil
}

func TestShortenHandler_Handle(t *testing.T) {
	mockService := new(ServiceMock)

	handler := NewShortenHandler(mockService, "http://example.com")

	type want struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "EmptyBody",
			body: "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"error\":\"EOF\"}",
			},
		},
		{
			name: "EmptyJson",
			body: "{}",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"error\":\"Key: 'ShortenRequest.URL' Error:Field validation for 'URL' failed on the 'required' tag\"}",
			},
		},

		{
			name: "NotJson",
			body: "<xml></xml>",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"error\":\"invalid character '\\u003c' looking for beginning of value\"}",
			},
		},

		{
			name: "EmptyUrl",
			body: "{\"url\": \"\"}",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"error\":\"Key: 'ShortenRequest.URL' Error:Field validation for 'URL' failed on the 'required' tag\"}",
			},
		},

		{
			name: "Success",
			body: "{\"url\": \"https://practicum.yandex.ru\"}",
			want: want{
				statusCode: http.StatusCreated,
				body:       "{\"result\":\"http://example.com/7CwAhsKq\"}",
			},
		},
		{
			name: "error",
			body: "{\"url\": \"error\"}",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"error\":\"error\"}",
			},
		},
		{
			name: "exists",
			body: "{\"url\": \"exists\"}",
			want: want{
				statusCode: http.StatusConflict,
				body:       "{\"result\":\"http://example.com/7CwAhsKq\"}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			context, _ := gin.CreateTestContext(response)
			context.Request = request

			handler.Handle(context)

			result := response.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			body, err := io.ReadAll(result.Body)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}
