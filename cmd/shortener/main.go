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
	"github.com/arturturundaev/shorturl/internal/app/handler/stats"
	"github.com/arturturundaev/shorturl/internal/app/handler/user"
	"github.com/arturturundaev/shorturl/internal/app/middleware"
	"github.com/arturturundaev/shorturl/internal/app/proto"
	"github.com/arturturundaev/shorturl/internal/app/repository/filestorage"
	"github.com/arturturundaev/shorturl/internal/app/repository/localstorage"
	pg "github.com/arturturundaev/shorturl/internal/app/repository/postgres"
	"github.com/arturturundaev/shorturl/internal/app/service"
	"github.com/arturturundaev/shorturl/internal/config"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Build variables
var (
	// BuildVersion - version
	BuildVersion string = "N/A"
	// BuildDate - date
	BuildDate string = "N/A"
	// BuildCommit - commit
	BuildCommit string = "N/A"
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

// UserUrlsCountStat
const UserUrlsCountStat = `/api/internal/stats`

func main() {
	startServe()
}

func startServe() {

	serverConfig, logger, jwtValidate, trustedSubnet, shortURLService, pingService, statService, repoRW := getResources()

	router := initRouter(shortURLService, serverConfig, pingService, statService, jwtValidate, trustedSubnet)

	logger.Info(fmt.Sprintf("Build version: %s\n", BuildVersion))
	logger.Info(fmt.Sprintf("Build date: %s\n", BuildDate))
	logger.Info(fmt.Sprintf("Build commit: %s\n", BuildCommit))

	if serverConfig.GRPCAddr != "" {

		// определяем порт для сервера
		listen, err := net.Listen("tcp", serverConfig.GRPCAddr)
		if err != nil {
			log.Fatal(err)
		}

		// создаём gRPC-сервер без зарегистрированной службы
		s := grpc.NewServer()
		// регистрируем сервис
		proto.RegisterServiceServer(s, proto.NewGPRCServer(pingService, statService))

		fmt.Println("Сервер gRPC начал работу " + serverConfig.GRPCAddr)
		// получаем запрос gRPC
		if err := s.Serve(listen); err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()
	} else {

		logger.Info("server start on port: " + serverConfig.AddressStart)
		server := &http.Server{
			Addr:    serverConfig.AddressStart,
			Handler: router,
		}

		go func() {
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			<-signalChan
			repoRW.SaveToFile("/tmp/save.txt")
			os.Exit(0)
		}()

		var errServer error
		if serverConfig.HTTPS.Enabled {
			errServer = server.ListenAndServeTLS(serverConfig.HTTPS.SSLPemPath, serverConfig.HTTPS.SSLKeyPath)
		} else {
			errServer = server.ListenAndServe()
		}
		if errServer != nil {
			logger.Fatal(errServer.Error())
		}
	}
}

func initRouter(shortURLService *service.ShortURLService, serverConfig *config.Config, pingService *service.PingService, statService *service.StatService, jwtValidate *middleware.JWTValidator, trustedSubnet *middleware.TrustedSubnetMiddleware) *gin.Engine {

	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	handlerFind := find.NewFindHandler(shortURLService)
	handlerSave := save.NewSaveHandler(shortURLService, serverConfig.BaseShort)
	handlerShortener := shorten.NewShortenHandler(shortURLService, serverConfig.BaseShort)
	handlerPing := ping.NewPingHandler(pingService)
	handlerButch := batch.NewButchHandler(shortURLService, serverConfig.BaseShort)
	handlerFindByUser := user.NewURLFindByUserHandler(shortURLService, serverConfig.BaseShort)
	handlerDelete := deleteUrl.NewDeleteHandler(shortURLService)
	handlerUsersUrlsCount := stats.NewUrlsAndUsersStatHandler(statService)

	router.GET(Ping, handlerPing.Handle)
	router.POST(SaveFullURL, jwtValidate.Handle, handlerSave.Handle)
	router.GET(GetFullURL, handlerFind.Handle)
	router.POST(SaveFullURL2, jwtValidate.Handle, handlerShortener.Handle)
	router.POST(SaveBatch, handlerButch.Handle)
	router.GET(URLByUser, jwtValidate.Handle, handlerFindByUser.Handle)
	router.DELETE(DeleteByUrls, jwtValidate.Handle, handlerDelete.Handle)
	router.GET(UserUrlsCountStat, trustedSubnet.Handle, handlerUsersUrlsCount.Handle)

	pprof.Register(router, "dev/pprof")

	return router
}

func getResources() (*config.Config, *zap.Logger, *middleware.JWTValidator, *middleware.TrustedSubnetMiddleware, *service.ShortURLService, *service.PingService, *service.StatService, service.RepositoryWriter) {

	serverConfig := config.NewConfig()

	logger, err := addLogger(nil, serverConfig.FullLog)

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

	jwtValidate := middleware.NewJWTValidator(serverConfig.AddressStart)
	trustedSubnet := middleware.NewTrustedSubnetMiddleware(serverConfig.TrustedSubnetFinal)

	shortURLService := service.NewShortURLService(repositoryRead, repositoryWrite)

	pingService := service.NewPingService(repositoryRead)

	statService := service.NewStatService(repositoryRead)

	return serverConfig, logger, jwtValidate, trustedSubnet, shortURLService, pingService, statService, repositoryWrite
}

func addLogger(r *gin.Engine, fullLogger bool) (*zap.Logger, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

	if fullLogger && r != nil {
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
		database, err := sqlx.Open("postgres", serverConfig.DatabaseURL)
		if err != nil {
			return nil, nil
		}
		repository, errPingRepo := pg.NewPostgresRepository(database)
		if errPingRepo != nil {
			logger.Error(errPingRepo.Error())
		}

		return repository, repository
	}

	if serverConfig.StorageType == config.StorageTypeFile {
		repositoryWrite, errStorageWrite := filestorage.NewFileStorageRepositoryWrite(serverConfig.FileStorage)
		if errStorageWrite != nil {
			logger.Error(errStorageWrite.Error())
		}

		repositoryRead, errStorageRead := filestorage.NewFileStorageRepositoryRead(serverConfig.FileStorage)
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
