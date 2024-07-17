package main

import (
	"flag"
	"github.com/arturturundaev/shorturl/internal/app/handler"
	"github.com/arturturundaev/shorturl/internal/app/repository/localstorage"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/arturturundaev/shorturl/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {

	serverConfig := config.NewConfig(os.Getenv("SERVER_ADDRESS"), os.Getenv("BASE_URL"))

	flag.Var(&serverConfig.AddressStart, "a", "start url and port")
	flag.Var(&serverConfig.BaseShort, "b", "url redirect")
	flag.Parse()

	repository := localstorage.NewLocalStorageRepository()
	shortUrlService := service.NewShortURLService(repository)
	handlerFind := handler.NewFindHandler(shortUrlService)
	handlerSave := handler.NewSaveHandler(shortUrlService, serverConfig.BaseShort.URL)

	router := gin.Default()

	router.POST(`/`, handlerSave.Handle)
	router.GET(`/:short`, handlerFind.Handle)

	err := http.ListenAndServe(serverConfig.AddressStart.String(), router)
	if err != nil {
		panic(err)
	}

	http.Server.Shutdown()
}
