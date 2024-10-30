package main

import (
	"bytes"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	deleteUrl "github.com/arturturundaev/shorturl/internal/app/handler/delete"
	"github.com/arturturundaev/shorturl/internal/app/handler/find"
	"github.com/arturturundaev/shorturl/internal/app/handler/ping"
	"github.com/arturturundaev/shorturl/internal/app/handler/save"
	"github.com/arturturundaev/shorturl/internal/app/handler/shorten"
	"github.com/arturturundaev/shorturl/internal/app/handler/user"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/arturturundaev/shorturl/internal/app/repository/filestorage"
	"github.com/arturturundaev/shorturl/internal/app/repository/localstorage"
	pg "github.com/arturturundaev/shorturl/internal/app/repository/postgres"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/arturturundaev/shorturl/internal/config"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	_ "github.com/gin-contrib/pprof"
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
	"path/filepath"
	"time"
)

// -d=postgres://postgres:postgres@localhost:5432/shorturl?sslmode=disable
// SaveFullURL создание
const SaveFullURL = `/`

// GetFullURL получение
const GetFullURL = `/:short`

// SaveFullURL2 v2
const SaveFullURL2 = `/api/shorten`

// SaveBatch массовое сохранение
const SaveBatch = `/api/shorten/batch`

// Ping пинг
const Ping = `/ping`

// URLByUser получение по пользователю
const URLByUser = `/api/user/urls`

// DeleteByUrls удаление ссылок
const DeleteByUrls = `/api/user/urls`

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	serverConfig := config.NewConfig()

	router := gin.Default()

	logger, err := addLogger(router, serverConfig.FullLog)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(5)
	}

	repositoryRead, repositoryWrite := getRepository(serverConfig, logger)

	if repositoryRead == nil {
		logger.Error("ошибка инициализации репозитория на чтение. Тип репозитория: " + config.StorageTypeDB)
	}

	if repositoryWrite == nil {
		logger.Error("ошибка инициализации репозитория на запись. Тип репозитория: " + config.StorageTypeDB)
	}

	if serverConfig.StorageType == config.StorageTypeDB {
		absPath, errPathMigration := filepath.Abs(".")
		if errPathMigration != nil {
			logger.Error("ошибка определения директории для миграций!")
		} else {
			initMigrations("file:////"+absPath+"/internal/app/repository/postgres/migration", repositoryRead.GetDB())
		}
	}

	jwtValidate := middleware.NewJWTValidator(serverConfig.AddressStart.URL)

	shortURLService := service.NewShortURLService(repositoryRead, repositoryWrite)

	pingService := service.NewPingService(repositoryRead)

	handlerFind := find.NewFindHandler(shortURLService)
	handlerSave := save.NewSaveHandler(shortURLService, serverConfig.BaseShort.URL)
	handlerSave2 := shorten.NewShortenHandler(shortURLService, serverConfig.BaseShort.URL)
	handlerPing := ping.NewPingHandler(pingService)
	handlerButch := batch.NewButchHandler(shortURLService, serverConfig.BaseShort.URL)
	handlerFindByUser := user.NewURLFindByUserHandler(shortURLService, serverConfig.BaseShort.URL)
	handlerDelete := deleteUrl.NewDeleteHandler(shortURLService)

	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	router.GET(Ping, handlerPing.Handle)
	router.POST(SaveFullURL, jwtValidate.Handle, handlerSave.Handle)
	router.GET(GetFullURL, handlerFind.Handle)
	router.POST(SaveFullURL2, jwtValidate.Handle, handlerSave2.Handle)
	router.POST(SaveBatch, handlerButch.Handle)
	router.GET(URLByUser, jwtValidate.Handle, handlerFindByUser.Handle)
	router.DELETE(DeleteByUrls, jwtValidate.Handle, handlerDelete.Handle)

	pprof.Register(router, "dev/pprof")

	logger.Info("server start on port: " + serverConfig.AddressStart.String())

	errServer := http.ListenAndServe(serverConfig.AddressStart.String(), router)
	if errServer != nil {
		logger.Fatal(errServer.Error())
	}
}

func addLogger(r *gin.Engine, fullLogger bool) (*zap.Logger, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

	if fullLogger {
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

	return logger, nil
}

func initMigrations(migrationPath string, DB *sqlx.DB) {
	driver, err := postgres.WithInstance(DB.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
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

func getRepository(serverConfig *config.Config, logger *zap.Logger) (service.RepositoryReader, service.RepositoryWriter) {

	if serverConfig.StorageType == config.StorageTypeDB {
		repository, errPingRepo := pg.NewPostgresRepository(serverConfig.DatabaseURL.URL)
		if errPingRepo != nil {
			logger.Error(errPingRepo.Error())
		}

		return repository, repository
	}

	if serverConfig.StorageType == config.StorageTypeFile {
		repositoryWrite, errStorageWrite := filestorage.NewFileStorageRepositoryWrite(serverConfig.FileStorage.Path)
		if errStorageWrite != nil {
			logger.Error(errStorageWrite.Error())
		}

		repositoryRead, errStorageRead := filestorage.NewFileStorageRepositoryRead(serverConfig.FileStorage.Path)
		if errStorageRead != nil {
			logger.Error(errStorageRead.Error())
		}

		return repositoryRead, repositoryWrite
	}

	if serverConfig.StorageType == config.StorageTypeMemory {
		repositoryWrite := localstorage.NewLocalStorageRepository()

		return repositoryWrite, repositoryWrite
	}

	return nil, nil
}
