package save

import (
	"errors"
	"io"
	"net/http"

	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
)

// SaveHandler сервис
type SaveHandler struct {
	service *service.ShortURLService
	baseURL string
}

// NewSaveHandler конструктор
func NewSaveHandler(service *service.ShortURLService, baseURL string) *SaveHandler {
	return &SaveHandler{service: service, baseURL: baseURL}
}

// Handle обработчик сохранения
func (hndlr *SaveHandler) Handle(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	data, err := hndlr.service.Save(ctx, string(b))

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
