package handler

import (
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockSaveRepository struct{}

func (repository *MockSaveRepository) FindByShortUrl(shortUrl string) (*entity.ShortUrlEntity, error) {

	return nil, nil
}

func (repository *MockSaveRepository) Save(shortUrl string, url string) error {

	if url == "repositoryError" {
		return fmt.Errorf("Error on insert row")
	}
	return nil
}

func TestSaveHandler_Handle(t *testing.T) {

	mockRepository := new(MockSaveRepository)

	handler := NewSaveHandler(service.NewShortUrlService(mockRepository), "http://example.com")

	type want struct {
		statusCode int
		body       string
		location   string
	}

	tests := []struct {
		name        string
		method      string
		contentType string
		request     string
		body        string
		want        want
	}{
		{
			name:    "Repository Error",
			method:  http.MethodPost,
			request: "/",
			body:    "repositoryError",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "Error on insert row",
			},
		},
		{
			name:    "Success",
			method:  http.MethodPost,
			request: "/",
			body:    "http://ya.ru",
			want: want{
				statusCode: http.StatusCreated,
				body:       "http://example.com/nnF0wba_",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			context, _ := gin.CreateTestContext(response)
			context.AddParam("short", strings.TrimLeft(tt.request, "/"))
			context.Request = request

			handler.Handle(context)

			result := response.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			body, err := io.ReadAll(result.Body)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}
