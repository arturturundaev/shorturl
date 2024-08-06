package shorten

import (
	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

func NewShortenRequest(context *gin.Context) (*ShortenRequest, error) {
	dto := &ShortenRequest{}

	if err := context.BindJSON(dto); err != nil {
		return dto, err
	}

	return dto, nil
}
