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
	dto, err := NewShortenRequest(ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dto.URL == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty URL"})
		return
	}

	data, err2 := h.service.Save(dto.URL)

	if err2 != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	response := ShortenResponse{URL: fmt.Sprintf("%s/%s", h.baseURL, data.ShortURL)}

	ctx.Writer.Header().Set("Accept-Encoding", "gzip")
	ctx.Writer.Header().Set("Content-Encoding", "gzip")
	ctx.Writer.Header().Set("Content-Type", "application/json")
	bt, _ := json.Marshal(response)
	ctx.Data(http.StatusOK, "gzip", bt)
}
