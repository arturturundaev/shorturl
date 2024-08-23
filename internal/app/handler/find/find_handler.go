package find

import (
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FindHandler struct {
	service *service.ShortURLService
}

func NewFindHandler(service *service.ShortURLService) *FindHandler {
	return &FindHandler{service: service}
}

func (hndlr *FindHandler) Handle(ctx *gin.Context) {

	data, err := hndlr.service.FindByShortURL(ctx.Param("short"))

	if err != nil || data == nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	if data.IsDeleted == true {
		ctx.Status(http.StatusGone)
		ctx.Abort()
	}

	ctx.Redirect(http.StatusTemporaryRedirect, data.URL)
}
