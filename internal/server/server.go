package server

import (
	"errors"
	"fmt"
	"log/slog"

	urlbookpkg "github.com/artm2000/urlbook/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Config struct {
	Host string
	Port string
}

type API struct{}

func (api *API) healthStatus(c *fiber.Ctx) error {
	return c.JSON(urlbookpkg.FormatResponse(c, urlbookpkg.Data{
		Data: map[string]interface{}{
			"server_status": "fine",
		},
	}))
}

func Run(config Config) func() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		// ReduceMemoryUsage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// check that if error was an fiber NewError and got status code,
			// specify that in error handler
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

			return c.Status(code).JSON(urlbookpkg.FormatResponse(c, urlbookpkg.Data{
				Message: err.Error(),
				IsError: true,
			}))
		},
	})

	/**
	 * General configuration
	 */
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${pid}] '${ip}:${port}' ${status} - ${method} ${path}\n",
	}))
	app.Use(requestid.New(requestid.Config{
		Next: func(c *fiber.Ctx) bool {
			trackId := c.Get(fiber.HeaderXRequestID)
			if trackId != "" {
				c.Set(fiber.HeaderXRequestID, trackId)
				return true
			}
			return false
		},
	}))
	app.Use(helmet.New())

	api := API{}

	app.Get("/v1/health", api.healthStatus)

	app.Hooks().OnListen(func(ld fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}
		scheme := "http"
		if ld.TLS {
			scheme = "https"
		}
		slog.Info(fmt.Sprintf("server start listening on '%s'", scheme+"://"+ld.Host+":"+ld.Port))
		return nil
	})

	go func() {
		if err := app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port)); err != nil {
			slog.Error(err.Error())
		}
	}()

	return func() {
		app.Shutdown()
	}
}
