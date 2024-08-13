package batch

import (
	"encoding/json"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ButchURLServiceInterface interface {
	Batch(request *[]ButchRequest) (*[]entity.ShortURLEntity, error)
}

type ButchHandler struct {
	service ButchURLServiceInterface
	baseURL string
}

func NewButchHandler(service ButchURLServiceInterface, baseURL string) *ButchHandler {
	return &ButchHandler{service: service, baseURL: baseURL}
}

func (h *ButchHandler) Handle(ctx *gin.Context) {
	var response []ButchResponse
	request, err := NewButchRequest(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models, err := h.service.Batch(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, model := range *models {
		response = append(response, ButchResponse{CorrelationId: model.CorrelationId, ShortURL: model.ShortURL})
	}

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