package handler

import (
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FindHandler struct {
	service *service.ShortUrlService
}

func NewFindHandler(service *service.ShortUrlService) *FindHandler {
	return &FindHandler{service: service}
}

func (hndlr *FindHandler) Handle(ctx *gin.Context) {

	data, err := hndlr.service.FindByShortUrl(ctx.Param("short"))

	if err != nil || data == nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, data.Url)

	return
}
