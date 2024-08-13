package batch

import "github.com/gin-gonic/gin"

type ButchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

func NewButchRequest(context *gin.Context) (*[]ButchRequest, error) {
	dto := &[]ButchRequest{}

	if err := context.BindJSON(dto); err != nil {
		return dto, err
	}

	return dto, nil
}
