package handler

import (
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockFindRepository struct{}

func (repository *MockFindRepository) FindByShortUrl(shortUrl string) (*entity.ShortUrlEntity, error) {

	if shortUrl == "repositoryError" {
		return nil, fmt.Errorf("Row not found by short url: %s", shortUrl)
	}

	if shortUrl == "find" {
		return &entity.ShortUrlEntity{Url: "findFull", ShortUrl: "find"}, nil
	}

	return nil, nil
}

func (repository *MockFindRepository) Save(shortUrl string, url string) error {

	return nil
}

func TestFindHandler_Handle(t *testing.T) {

	mockRepository := new(MockFindRepository)

	handler := NewFindHandler(service.NewShortUrlService(mockRepository))

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
			name:    "Invalid METHOD type",
			method:  http.MethodPost,
			request: "/qwerty",
			body:    "",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "Only GET requests are allowed!\n",
				location:   "",
			},
		},
		{
			name:    "Repository Error",
			method:  http.MethodGet,
			request: "/repositoryError",
			body:    "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "Row not found by short url: repositoryError\n",
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

			handler.Handle(response, request)

			result := response.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			body, err := io.ReadAll(result.Body)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.body, string(body))
			assert.Equal(t, result.Header.Get("Location"), tt.want.location)
		})
	}
}
