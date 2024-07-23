package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/handler/find"
	"github.com/arturturundaev/shorturl/internal/app/handler/save"
	"github.com/arturturundaev/shorturl/internal/app/handler/shorten"
	"github.com/arturturundaev/shorturl/internal/app/repository/filestorage"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/arturturundaev/shorturl/internal/config"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"os"
	"time"
)

const SAVE_FULL_URL = `/`
const GET_FULL_URL = `/:short`
const SAVE_FULL_URL_2 = `/api/shorten`

func main() {

	serverConfig := config.NewConfig(os.Getenv("SERVER_ADDRESS"), os.Getenv("BASE_URL"))

	flag.Var(&serverConfig.AddressStart, "a", "start url and port")
	flag.Var(&serverConfig.BaseShort, "b", "url redirect")
	flag.Parse()

	repositoryWrite, err := filestorage.NewFileStorageRepositoryWrite("db.txt")
	if err != nil {
		panic(err)
	}

	repositoryRead, err2 := filestorage.NewFileStorageRepositoryRead("db.txt")
	if err2 != nil {
		panic(err)
	}

	shortURLService := service.NewShortURLService(repositoryRead, repositoryWrite)
	handlerFind := find.NewFindHandler(shortURLService)
	handlerSave := save.NewSaveHandler(shortURLService, serverConfig.BaseShort.URL)
	handlerSave2 := shorten.NewShortenHandler(shortURLService, serverConfig.BaseShort.URL)

	router := gin.Default()

	addLogger(router)

	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	router.POST(SAVE_FULL_URL, handlerSave.Handle)
	router.GET(GET_FULL_URL, handlerFind.Handle)
	router.POST(SAVE_FULL_URL_2, handlerSave2.Handle)

	fmt.Println(serverConfig.AddressStart.String())
	err3 := http.ListenAndServe(serverConfig.AddressStart.String(), router)
	if err3 != nil {
		panic(err)
	}
}

func addLogger(r *gin.Engine) {
	logger, _ := zap.NewProduction()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			var fields []zapcore.Field
			// log request ID
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			// log request body
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
			fields = append(fields, zap.String("body", string(body)))

			return fields
		}),
	}))
}
