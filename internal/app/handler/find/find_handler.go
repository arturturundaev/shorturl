package find

import (
	"fmt"
	"net/http"

	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
)

// FindHandler сервис
type FindHandler struct {
	service *service.ShortURLService
}

// NewFindHandler конструктор
func NewFindHandler(service *service.ShortURLService) *FindHandler {
	return &FindHandler{service: service}
}

// Handle обработчик поиска
func (hndlr *FindHandler) Handle(ctx *gin.Context) {

	data, err := hndlr.service.FindByShortURL(ctx.Param("short"))

	if err != nil || data == nil {
		er := err
		if data == nil {
			er = fmt.Errorf("not find")
		}
		ctx.String(http.StatusBadRequest, "%s", er.Error())
		ctx.Abort()
		return
	}

	if data.IsDeleted {
		ctx.Status(http.StatusGone)
		ctx.Abort()
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, data.URL)
}
