package shorten

import (
	"encoding/json"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
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
	dto, errGenerateShortUrl := NewShortenRequest(ctx)

	if errGenerateShortUrl != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errGenerateShortUrl.Error()})
		return
	}

	if dto.URL == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty URL"})
		return
	}

	data, errRepository := h.service.Save(dto.URL)

	if errRepository != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errRepository.Error()})
		return
	}

	response := ShortenResponse{URL: fmt.Sprintf("%s/%s", h.baseURL, data.ShortURL)}
	bt, errMarshal := json.Marshal(response)

	if errMarshal != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errMarshal.Error()})
		return
	}

	ctx.Writer.Header().Set("Accept-Encoding", "gzip")
	ctx.Writer.Header().Set("Content-Encoding", "gzip")
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Data(http.StatusCreated, "gzip", bt)
}
