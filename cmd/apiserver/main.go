package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/artm2000/urlbook/internal/controller"
	"github.com/artm2000/urlbook/internal/core/common"
	"github.com/artm2000/urlbook/internal/core/server"
	"github.com/artm2000/urlbook/internal/core/service"
	"github.com/artm2000/urlbook/internal/infra/config"
	"github.com/artm2000/urlbook/internal/infra/repository"

	"github.com/artm2000/urlbook/pkg"
	"github.com/joho/godotenv"
)

/**
 * env keys
 */
const (
	SERVER_HOST_KEY      string = "SERVER_HOST"
	SERVER_PORT_KEY      string = "SERVER_PORT"
	DATABASE_HOST        string = "DATABASE_HOST"
	DATABASE_NAME        string = "DATABASE_NAME"
	DATABASE_USER        string = "DATABASE_USER"
	DATABASE_PASSWORD    string = "DATABASE_PASSWORD"
	DATABASE_PUBLIC_PORT string = "DATABASE_PUBLIC_PORT"
	MEMCACHED_ADDRESS    string = "MEMCACHED_ADDRESS"
	PUBLIC_ADDRESS       string = "PUBLIC_ADDRESS"
	CUSTOM_ENV_FILE      string = "ENV_FILE"
)

const CUSTOM_ENV_FILE_FLAG string = "envfile"

var envFileFlag = flag.String(CUSTOM_ENV_FILE_FLAG, "", "define server env file to read its required values")

func main() {
	/**
	 * Prepare logger
	 */
	slog.SetDefault(common.NewLogger(slog.LevelDebug))

	/**
	 * Try to load env file, from flag or environment variables
	 */
	envFile, envFileExist := os.LookupEnv(CUSTOM_ENV_FILE)
	if !envFileExist {
		flag.Parse()
		envFile = *envFileFlag
	}
	if err := godotenv.Load(envFile); err != nil {
		slog.Error(err.Error())
	}

	dbConnection := config.NewMysqlDBConn(
		getValueFromEnv(DATABASE_USER),
		getValueFromEnv(DATABASE_PASSWORD),
		getValueFromEnv(DATABASE_HOST),
		getValueFromEnv(DATABASE_PUBLIC_PORT),
		getValueFromEnv(DATABASE_NAME),
	)
	memcachedClient := config.NewMemcachedClient(getValueFromEnv(MEMCACHED_ADDRESS))

	memcachedRepository := repository.NewMemcachedRepository(memcachedClient)
	urlRepository := repository.NewUrlRepository(dbConnection)
	urlShortenerService := service.NewUrlShortener(
		urlRepository,
		memcachedRepository,
		tryGetValueFromEnv(PUBLIC_ADDRESS, ""),
	)

	h := server.NewHttpServer(config.HttpServer{
		Host: tryGetValueFromEnv(SERVER_HOST_KEY, "127.0.0.1"),
		Port: tryGetValueFromEnv(SERVER_PORT_KEY, "3000"),
	})
	h.Start()
	h.RegisterControllers(
		controller.NewUrlShortener(
			urlShortenerService,
		),
		controller.NewUrlRedirect(
			urlShortenerService,
		),
	)

	pkg.OnInterrupt(func() { h.Stop(true) })
}

func tryGetValueFromEnv(key, defaultValue string) string {
	value, valueExist := os.LookupEnv(key)
	if !valueExist {
		return defaultValue
	}
	return value
}

func getValueFromEnv(key string) string {
	value, valueExist := os.LookupEnv(key)
	if !valueExist {
		panic(fmt.Sprintf("fail to get '%s' from environment variables", key))
	}
	return value
}
