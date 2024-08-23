package save

import (
	"context"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockReadRepository struct{}

func (repository *MockReadRepository) GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *MockReadRepository) Ping(ctx context.Context) error {
	return nil
}

func (repository *MockReadRepository) GetDB() *sqlx.DB {
	return nil
}

type MockWriteRepository struct{}

func (repository *MockWriteRepository) Delete(shortURLs []string, addedUserID string) error {
	//TODO implement me
	panic("implement me")
}

func (repository *MockWriteRepository) Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	return nil, nil
}

func (repository *MockWriteRepository) GetDB() *sqlx.DB {
	return nil
}

func (repository *MockReadRepository) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {

	return nil, nil
}

func (repository *MockWriteRepository) Save(shortURL string, url string, addedUserID string) error {

	if url == "repositoryError" {
		return fmt.Errorf("Error on insert row")
	}
	return nil
}

func TestSaveHandler_Handle(t *testing.T) {

	mockReadRepository := new(MockReadRepository)
	mockWriteRepository := new(MockWriteRepository)

	handler := NewSaveHandler(service.NewShortURLService(mockReadRepository, mockWriteRepository), "http://example.com")

	type want struct {
		statusCode int
		body       string
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
			context.Set(middleware.UserIDProperty, "1")
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
