package main

import (
	"github.com/arturturundaev/shorturl/internal/app/handler"
	"github.com/arturturundaev/shorturl/internal/app/repository/local_storage"
	service2 "github.com/arturturundaev/shorturl/internal/app/service"
	"net/http"
)

var LocalStorage = make(map[string]string)

func main() {

	repository := local_storage.NewLocalStorageRepository()
	service := service2.NewShortUrlService(repository)
	handlerFind := handler.NewFindHandler(service)
	handlerSave := handler.NewSaveHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handlerSave.Handle)
	mux.HandleFunc(`/{id}`, handlerFind.Handle)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
