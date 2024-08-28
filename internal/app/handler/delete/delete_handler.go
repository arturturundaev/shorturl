package delete

import (
	"encoding/json"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type DeleteHandler struct {
	service *service.ShortURLService
}

func NewDeleteHandler(service *service.ShortURLService) *DeleteHandler {
	return &DeleteHandler{service: service}
}

func (h *DeleteHandler) Handle(ctx *gin.Context) {
	var data []string
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	addedUserID, exists := ctx.Get(middleware.UserIDProperty)

	if !exists {
		ctx.Status(http.StatusUnauthorized)
	}

	go h.service.Delete(data, addedUserID.(string))

	ctx.Status(http.StatusAccepted)
}
