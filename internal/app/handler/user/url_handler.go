package user

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type URLServiceFinder interface {
	GetUrlsByUserID(userID string) ([]entity.ShortURLEntity, error)
}

type URLFindByUserHandler struct {
	service URLServiceFinder
	baseURL string
}

func NewURLFindByUserHandler(service URLServiceFinder, baseURL string) *URLFindByUserHandler {
	return &URLFindByUserHandler{service: service, baseURL: baseURL}
}

func (handler *URLFindByUserHandler) Handle(ctx *gin.Context) {
	addedUserID, _ := ctx.Get(middleware.UserIDProperty)

	data, err := handler.service.GetUrlsByUserID(addedUserID.(string))

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	if len(data) == 0 {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	var response []URLListItemResponse

	for _, url := range data {
		response = append(response, NewURLResponse(handler.baseURL, url.ShortURL, url.URL))
	}
	ctx.JSON(http.StatusOK, response)
}
