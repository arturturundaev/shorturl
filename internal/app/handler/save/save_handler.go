package save

import (
	"errors"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type SaveHandler struct {
	service *service.ShortURLService
	baseURL string
}

func NewSaveHandler(service *service.ShortURLService, baseURL string) *SaveHandler {
	return &SaveHandler{service: service, baseURL: baseURL}
}

func (hndlr *SaveHandler) Handle(ctx *gin.Context) {

	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	data, err := hndlr.service.Save(string(b))

	if errors.Is(err, service.ErrEntityExists) {
		ctx.Header("Content-type", "text/plain")
		ctx.String(http.StatusConflict, "%s/%s", hndlr.baseURL, data.ShortURL)
		return
	}

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	ctx.Header("Content-type", "text/plain")
	ctx.String(http.StatusCreated, "%s/%s", hndlr.baseURL, data.ShortURL)
}
