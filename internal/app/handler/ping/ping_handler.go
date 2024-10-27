package ping

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type pignger interface {
	Ping(ctx context.Context) error
}

type PingHandler struct {
	service pignger
}

func NewPingHandler(service pignger) *PingHandler {
	return &PingHandler{service: service}
}

func (h *PingHandler) Handle(ctx *gin.Context) {

	contxt, cancel := context.WithCancel(ctx)
	defer cancel()
	time.AfterFunc(1500*time.Millisecond, cancel)

	err := h.service.Ping(contxt)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, "")
}
