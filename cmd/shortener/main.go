package main

import (
	"github.com/arturturundaev/shorturl/internal/app/handler"
	"github.com/arturturundaev/shorturl/internal/app/repository/local_storage"
	service2 "github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var LocalStorage = make(map[string]string)

func main() {

	repository := local_storage.NewLocalStorageRepository()
	service := service2.NewShortUrlService(repository)
	handlerFind := handler.NewFindHandler(service)
	handlerSave := handler.NewSaveHandler(service)

	router := gin.Default()

	router.POST(`/`, handlerSave.Handle)
	router.GET(`/:short`, handlerFind.Handle)

	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
