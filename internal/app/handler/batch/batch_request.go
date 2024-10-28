package batch

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// ButchRequest dto
type ButchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// NewButchRequest конструктор
func NewButchRequest(context *gin.Context) ([]ButchRequest, error) {
	var dto []ButchRequest

	if err := context.BindJSON(&dto); err != nil {
		return dto, err
	}

	if len(dto) == 0 {
		return nil, fmt.Errorf("пустой входной массив")
	}

	return dto, nil
}
