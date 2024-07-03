package handler

import (
	"github.com/arturturundaev/shorturl/internal/app/service"
	"io"
	"net/http"
)

type SaveHandler struct {
	service *service.ShortUrlService
}

func NewSaveHandler(service *service.ShortUrlService) *SaveHandler {
	return &SaveHandler{service: service}
}

func (hndlr *SaveHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := hndlr.service.Save(string(b))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write([]byte("http://" + r.Host + "/" + data.ShortUrl)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "text/plain")

	return
}
