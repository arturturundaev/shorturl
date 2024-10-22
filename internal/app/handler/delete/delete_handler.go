package delete

import (
	"encoding/json"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type deleter interface {
	Delete(URLList []string, addedUserID string)
}

type DeleteHandler struct {
	service deleter
}

func NewDeleteHandler(service deleter) *DeleteHandler {
	return &DeleteHandler{service: service}
}

func (h *DeleteHandler) Handle(ctx *gin.Context) {
	var data []string
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	addedUserID, exists := ctx.Get(middleware.UserIDProperty)

	if !exists {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	go h.service.Delete(data, addedUserID.(string))

	ctx.Status(http.StatusAccepted)
}
