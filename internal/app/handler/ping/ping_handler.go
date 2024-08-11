package ping

import (
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PingHandler struct {
	service *service.PingService
}

func NewPingHandler(service *service.PingService) *PingHandler {
	return &PingHandler{service: service}
}

func (h *PingHandler) Handle(ctx *gin.Context) {
	err := h.service.Ping(ctx)

	if err != nil {
		fmt.Printf(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, "")
	return
}
