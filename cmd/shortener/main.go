package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/handler/find"
	"github.com/arturturundaev/shorturl/internal/app/handler/ping"
	"github.com/arturturundaev/shorturl/internal/app/handler/save"
	"github.com/arturturundaev/shorturl/internal/app/handler/shorten"
	pg "github.com/arturturundaev/shorturl/internal/app/repository/postgres"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/arturturundaev/shorturl/internal/config"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
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
const Ping = `/ping`

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	serverConfig := config.NewConfig(
		os.Getenv("SERVER_ADDRESS"),
		os.Getenv("BASE_URL"),
		os.Getenv("FILE_STORAGE_PATH"),
		os.Getenv("DATABASE_DSN"))

	flag.Var(&serverConfig.AddressStart, "a", "start url and port")
	flag.Var(&serverConfig.BaseShort, "b", "url redirect")
	flag.Var(&serverConfig.FileStorage, "f", "file storage path")
	flag.Var(&serverConfig.DatabaseURL, "d", "database storage path")
	flag.Parse()

	router := gin.Default()

	logger, err := addLogger(router)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(5)
	}

	/*repositoryWrite, errStorageWrite := filestorage.NewFileStorageRepositoryWrite(serverConfig.FileStorage.Path)
	if errStorageWrite != nil {
		logger.Error(errStorageWrite.Error())
	}

	repositoryRead, errStorageRead := filestorage.NewFileStorageRepositoryRead(serverConfig.FileStorage.Path)
	if errStorageRead != nil {
		logger.Error(errStorageRead.Error())
	}*/

	postgresRepository, errPingRepo := pg.NewPostgresRepository(serverConfig.DatabaseURL.URL)
	if errPingRepo != nil {
		logger.Error(errPingRepo.Error())
	}

	initMigrations("file:////home/a_turundaev/projects/yandex/shorturl/internal/app/repository/postgres/migration", postgresRepository.DB)

	shortURLService := service.NewShortURLService(postgresRepository, postgresRepository)
	pingService := service.NewPingService(postgresRepository)

	handlerFind := find.NewFindHandler(shortURLService)
	handlerSave := save.NewSaveHandler(shortURLService, serverConfig.BaseShort.URL)
	handlerSave2 := shorten.NewShortenHandler(shortURLService, serverConfig.BaseShort.URL)
	handlerPing := ping.NewPingHandler(pingService)

	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	router.GET(Ping, handlerPing.Handle)
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

func initMigrations(migrationPath string, DB *sqlx.DB) {
	driver, err := postgres.WithInstance(DB.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)

	if err != nil {
		log.Fatal(err)
	} else {
		errMigrate := m.Up()
		if errMigrate != nil && errMigrate.Error() != "no change" {
			log.Fatal(errMigrate)
		}
	}
}
