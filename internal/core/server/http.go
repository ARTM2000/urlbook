package server

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/artm2000/urlbook/internal/controller"
	"github.com/artm2000/urlbook/internal/infra/config"
	urlbookpkg "github.com/artm2000/urlbook/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type HttpServer interface {
	Start()
	Stop(bool)
	RegisterControllers(...controller.HttpController)
}

type httpServer struct {
	config config.HttpServer
	app    *fiber.App
}

func NewHttpServer(config config.HttpServer) HttpServer {
	return &httpServer{
		config: config,
		app: nil,
	}
}

func (hs *httpServer) Start() {
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
		if err := app.Listen(fmt.Sprintf("%s:%s", hs.config.Host, hs.config.Port)); err != nil {
			slog.Error(err.Error())
		}
	}()

	hs.app = app
}

func (hs *httpServer) Stop(force bool) {
	if hs.app == nil {
		panic("http app server not defined. first start the http server")
	}
	if force {
		hs.app.Shutdown()
		return
	}
	hs.app.ShutdownWithTimeout(time.Second * 30)
}

func (hs *httpServer) RegisterControllers(controllers ...controller.HttpController) {
	if hs.app == nil {
		panic("http app server not defined. first start the http server")
	}
	for _, c := range controllers {
		c.InitRoutes(hs.app.Group(c.GetPrefix()))
	}
}
