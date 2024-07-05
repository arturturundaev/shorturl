package main

import (
	"flag"
	"github.com/arturturundaev/shorturl/internal/app/handler"
	"github.com/arturturundaev/shorturl/internal/app/repository/localstorage"
	service2 "github.com/arturturundaev/shorturl/internal/app/service"
	config2 "github.com/arturturundaev/shorturl/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

var LocalStorage = make(map[string]string)

func main() {

	config := config2.NewConfig()
	flag.Var(&config.AddressStart, "a", "start url and port")
	flag.Var(&config.BaseShort, "b", "url redirect")
	flag.Parse()

	repository := localstorage.NewLocalStorageRepository()
	service := service2.NewShortURLService(repository)
	handlerFind := handler.NewFindHandler(service)
	handlerSave := handler.NewSaveHandler(service, config.BaseShort.URL)

	router := gin.Default()

	router.POST(`/`, handlerSave.Handle)
	router.GET(`/:short`, handlerFind.Handle)

	err := http.ListenAndServe(config.AddressStart.String(), router)
	if err != nil {
		panic(err)
	}
}
