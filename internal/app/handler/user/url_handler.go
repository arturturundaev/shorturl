package user

import (
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type URLServiceFinder interface {
	GetUrlsByUserId(userId string) ([]entity.ShortURLEntity, error)
}

type URLFindByUserHandler struct {
	service URLServiceFinder
	baseURL string
}

func NewURLFindByUserHandler(service URLServiceFinder, baseURL string) *URLFindByUserHandler {
	return &URLFindByUserHandler{service: service, baseURL: baseURL}
}

func (handler *URLFindByUserHandler) Handle(ctx *gin.Context) {
	addedUserId, _ := ctx.Get(middleware.USER_ID_PROPERTY)

	data, err := handler.service.GetUrlsByUserId(addedUserId.(string))

	if err != nil {
		ctx.String(http.StatusBadRequest, "%s", err.Error())
		ctx.Abort()
		return
	}

	if len(data) == 0 {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	var response []UrlListItemResponse

	for _, url := range data {
		response = append(response, NewUrlResponse(handler.baseURL, url.ShortURL, url.URL))
	}
	ctx.JSON(http.StatusOK, response)
}
