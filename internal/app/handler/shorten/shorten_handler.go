package shorten

import (
	"errors"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SaveURLInterface interface {
	Save(url string) (*entity.ShortURLEntity, error)
}

type ShortenHandler struct {
	service SaveURLInterface
	baseURL string
}

func NewShortenHandler(service SaveURLInterface, baseURL string) *ShortenHandler {
	return &ShortenHandler{service: service, baseURL: baseURL}
}

func (h *ShortenHandler) Handle(ctx *gin.Context) {
	dto, errGenerateShortURL := NewShortenRequest(ctx)

	if errGenerateShortURL != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errGenerateShortURL.Error()})
		return
	}

	fmt.Println("JSON ????????????" + dto.URL + "??????????")
	if dto.URL == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty URL"})
		return
	}

	data, errRepository := h.service.Save(dto.URL)

	if errRepository != nil && !errors.Is(errRepository, service.EntityExistsError) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errRepository.Error()})
		return
	}

	response := ShortenResponse{URL: fmt.Sprintf("%s/%s", h.baseURL, data.ShortURL)}

	status := http.StatusCreated

	if errors.Is(errRepository, service.EntityExistsError) {
		status = http.StatusConflict
	}

	ctx.JSON(status, response)
}
