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
	"log"
	"net/http"
	"os"
	"time"
)

const SaveFullURL = `/`
const GetFullURL = `/:short`
const SaveFullURL2 = `/api/shorten`

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	serverConfig := config.NewConfig(os.Getenv("SERVER_ADDRESS"), os.Getenv("BASE_URL"), os.Getenv("FILE_STORAGE_PATH"))

	flag.Var(&serverConfig.AddressStart, "a", "start url and port")
	flag.Var(&serverConfig.BaseShort, "b", "url redirect")
	flag.Var(&serverConfig.FileStorage, "f", "file storage path")
	flag.Parse()

	router := gin.Default()

	logger, err := addLogger(router)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(5)
	}

	repositoryWrite, errStorageWrite := filestorage.NewFileStorageRepositoryWrite(serverConfig.FileStorage.Path)
	if errStorageWrite != nil {
		logger.Fatal(errStorageWrite.Error())
	}

	repositoryRead, errStorageRead := filestorage.NewFileStorageRepositoryRead(serverConfig.FileStorage.Path)
	if errStorageRead != nil {
		logger.Fatal(errStorageRead.Error())
	}

	shortURLService := service.NewShortURLService(repositoryRead, repositoryWrite)
	handlerFind := find.NewFindHandler(shortURLService)
	handlerSave := save.NewSaveHandler(shortURLService, serverConfig.BaseShort.URL)
	handlerSave2 := shorten.NewShortenHandler(shortURLService, serverConfig.BaseShort.URL)

	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	router.POST(SaveFullURL, handlerSave.Handle)
	router.GET(GetFullURL, handlerFind.Handle)
	router.POST(SaveFullURL2, handlerSave2.Handle)

	fmt.Println(">>>>>>> " + serverConfig.AddressStart.String() + " <<<<<<<<<")
	errServer := http.ListenAndServe(serverConfig.AddressStart.String(), router)
	if errServer != nil {
		logger.Fatal(errServer.Error())
	}
}

func addLogger(r *gin.Engine) (*zap.Logger, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

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

	return logger, nil
}
