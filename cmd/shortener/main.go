package main

import (
	"flag"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/handler"
	"github.com/arturturundaev/shorturl/internal/app/repository/local_storage"
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

	repository := local_storage.NewLocalStorageRepository()
	service := service2.NewShortUrlService(repository)
	handlerFind := handler.NewFindHandler(service)
	handlerSave := handler.NewSaveHandler(service, config.BaseShort.Url)

	router := gin.Default()

	router.POST(`/`, handlerSave.Handle)
	router.GET(`/:short`, handlerFind.Handle)

	fmt.Println(">>>>>>>> SERVER:" + config.AddressStart.String() + "<<<<<<<<<<<<<<<<<<")
	err := http.ListenAndServe(config.AddressStart.String(), router)
	if err != nil {
		panic(err)
	}
}
