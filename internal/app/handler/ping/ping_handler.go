package ping

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type pignger interface {
	Ping(ctx context.Context) error
}

// PingHandler сервис
type PingHandler struct {
	service pignger
}

// NewPingHandler конструктор
func NewPingHandler(service pignger) *PingHandler {
	return &PingHandler{service: service}
}

// Handle обработчик поиска
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
