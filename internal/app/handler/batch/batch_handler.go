package batch

import (
	"encoding/json"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type serviceURLButcher interface {
	Batch(request []ButchRequest) ([]entity.ShortURLEntity, error)
}

type ButchHandler struct {
	service serviceURLButcher
	baseURL string
}

func NewButchHandler(service serviceURLButcher, baseURL string) *ButchHandler {
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

	for _, model := range models {
		response = append(response, ButchResponse{CorrelationID: model.CorrelationID, ShortURL: fmt.Sprintf("%s/%s", h.baseURL, model.ShortURL)})
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
