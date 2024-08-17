package find

import (
	"context"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
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

func (repository *MockReadRepository) Ping(ctx context.Context) error {
	return nil
}

func (repository *MockReadRepository) GetDB() *sqlx.DB {
	return nil
}

type MockWriteRepository struct{}

func (repository *MockWriteRepository) Batch(request []batch.ButchRequest) ([]entity.ShortURLEntity, error) {
	return nil, nil
}

func (repository *MockWriteRepository) GetDB() *sqlx.DB {
	return nil
}

func (repository *MockReadRepository) FindByShortURL(shortURL string) (*entity.ShortURLEntity, error) {

	if shortURL == "repositoryError" {
		return nil, fmt.Errorf("Row not found by short url: %s", shortURL)
	}

	if shortURL == "find" {
		return &entity.ShortURLEntity{URL: "findFull", ShortURL: "find"}, nil
	}

	return nil, nil
}

func (repository *MockWriteRepository) Save(shortURL string, url string) error {

	return nil
}

func TestFindHandler_Handle(t *testing.T) {

	mockReadRepository := new(MockReadRepository)
	mockWriteRepository := new(MockWriteRepository)

	handler := NewFindHandler(service.NewShortURLService(mockReadRepository, mockWriteRepository))

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
			method:  http.MethodGet,
			request: "/repositoryError",
			body:    "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "Row not found by short url: repositoryError",
				location:   "",
			},
		},
		{
			name:    "Find",
			method:  http.MethodGet,
			request: "/find",
			body:    "",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				body:       "<a href=\"/findFull\">Temporary Redirect</a>.\n\n",
				location:   "/findFull",
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

			assert.Equal(t, tt.want.statusCode, response.Code)

			body, err := io.ReadAll(response.Body)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.body, string(body))
			assert.Equal(t, response.Header().Get("Location"), tt.want.location)
		})
	}
}
