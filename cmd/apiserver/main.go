package main

import (
	"flag"
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
const SERVER_HOST_KEY string = "SERVER_HOST"
const SERVER_PORT_KEY string = "SERVER_PORT"
const PUBLIC_ADDRESS = "PUBLIC_ADDRESS"
const CUSTOM_ENV_FILE string = "ENV_FILE"

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

	h := server.NewHttpServer(config.HttpServer{
		Host: getValueFromEnv(SERVER_HOST_KEY, "127.0.0.1"),
		Port: getValueFromEnv(SERVER_PORT_KEY, "3000"),
	})
	h.Start()
	h.RegisterControllers(
		controller.NewUrlShortener(
			service.NewUrlShortener(
				repository.NewUrlRepository(
					config.NewMysqlDBConfig("username", "password", "host", "port", "dbname"),
				),
				getValueFromEnv(PUBLIC_ADDRESS, ""),
			),
		),
	)

	pkg.OnInterrupt(func() { h.Stop(true) })
}

func getValueFromEnv(key, defaultValue string) string {
	value, valueExist := os.LookupEnv(key)
	if !valueExist {
		return defaultValue
	}
	return value
}
