package handler

import (
	"github.com/arturturundaev/shorturl/internal/app/service"
	"net/http"
	"strings"
)

type FindHandler struct {
	service *service.ShortUrlService
}

func NewFindHandler(service *service.ShortUrlService) *FindHandler {
	return &FindHandler{service: service}
}

func (hndlr *FindHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	data, err := hndlr.service.FindByShortUrl(strings.TrimLeft(r.RequestURI, "/"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, data.Url, http.StatusTemporaryRedirect)

	return
}
