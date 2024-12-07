package stats

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UrlsAndUsersStatHandler struct {
	service UrlsAndUsersStatService
}

type UrlsAndUsersStatService interface {
	GetUrlsAndUsersStat() (int32, int32)
}

func NewUrlsAndUsersStatHandler(service UrlsAndUsersStatService) *UrlsAndUsersStatHandler {
	return &UrlsAndUsersStatHandler{service: service}
}

func (h *UrlsAndUsersStatHandler) Handle(ctx *gin.Context) {
	urlsCount, usersCount := h.service.GetUrlsAndUsersStat()

	type response struct {
		Urls  int32 `json:"urls"`
		Users int32 `json:"users"`
	}

	ctx.JSON(http.StatusOK, response{
		Urls:  urlsCount,
		Users: usersCount,
	})
}
