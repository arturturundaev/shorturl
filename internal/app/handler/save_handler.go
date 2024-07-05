package handler

import (
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type SaveHandler struct {
	service *service.ShortUrlService
	baseUrl string
}

func NewSaveHandler(service *service.ShortUrlService, baseUrl string) *SaveHandler {
	return &SaveHandler{service: service, baseUrl: baseUrl}
}

func (hndlr *SaveHandler) Handle(ctx *gin.Context) {

	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	data, err := hndlr.service.Save(string(b))

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	ctx.Header("Content-type", "text/plain")
	ctx.String(http.StatusCreated, "%s/%s", hndlr.baseUrl, data.ShortUrl)
}
