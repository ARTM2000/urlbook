package controller

import (
	"context"
	"log/slog"

	"github.com/artm2000/urlbook/internal/core/model/request"
	"github.com/artm2000/urlbook/internal/core/port/service"
	"github.com/gofiber/fiber/v2"
)

type urlShortener struct {
	urlShortenerService service.UrlShortener
}

func NewUrlShortener(urlShortenerService service.UrlShortener) HttpController {
	return &urlShortener{
		urlShortenerService,
	}
}

func (usc *urlShortener) InitRoutes(api fiber.Router) {
	api.Post("/submit/", usc.submitUrl)
}

func (usc *urlShortener) GetPrefix() string { return "/v1/url" }

func (usc *urlShortener) submitUrl(c *fiber.Ctx) error {
	var body request.SubmitUrl
	body.TrackId = c.GetRespHeader(fiber.HeaderXRequestID)
	if err := c.BodyParser(&body); err != nil {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"error parsing request body",
			slog.Any("error", err.Error()),
		)
		return fiber.ErrBadRequest
	}

	result := usc.urlShortenerService.ShortUrl(&body)
	slog.LogAttrs(
		context.Background(),
		slog.LevelDebug,
		"submit-url response",
		slog.Any("result", result),
	)
	return c.Status(result.StatusCode).JSON(result)
}
