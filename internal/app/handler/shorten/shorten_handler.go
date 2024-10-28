package shorten

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
)

// Handle интерфейс для сохранения
type URLSaver interface {
	Save(ctx *gin.Context, url string) (*entity.ShortURLEntity, error)
}

// Handle сервис
type ShortenHandler struct {
	service URLSaver
	baseURL string
}

// Handle конструктор
func NewShortenHandler(service URLSaver, baseURL string) *ShortenHandler {
	return &ShortenHandler{service: service, baseURL: baseURL}
}

// Handle обаботчик
func (h *ShortenHandler) Handle(ctx *gin.Context) {
	dto, errGenerateShortURL := NewShortenRequest(ctx)

	if errGenerateShortURL != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errGenerateShortURL.Error()})
		return
	}

	if dto.URL == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty URL"})
		return
	}

	data, errRepository := h.service.Save(ctx, dto.URL)

	if errRepository != nil && !errors.Is(errRepository, service.ErrEntityExists) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errRepository.Error()})
		return
	}

	response := ShortenResponse{URL: fmt.Sprintf("%s/%s", h.baseURL, data.ShortURL)}

	status := http.StatusCreated

	if errors.Is(errRepository, service.ErrEntityExists) {
		status = http.StatusConflict
	}

	ctx.JSON(status, response)
}
