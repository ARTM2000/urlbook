package controller

import (
	"context"
	"log/slog"

	"github.com/artm2000/urlbook/internal/core/common"
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
	api.Post("/submit/custom/", usc.submitUrlByCustomPhrase)
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
		return fiber.ErrUnprocessableEntity
	}

	if errs, ok := common.GetValidator().ValidateStruct(&body); !ok {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"validation error",
			slog.String("error", errs[0].Message),
		)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	result := usc.urlShortenerService.ShortUrl(&body)
	return c.Status(result.StatusCode).JSON(result)
}

func (usc *urlShortener) submitUrlByCustomPhrase(c *fiber.Ctx) error {
	var body request.SubmitUrlByCustomPhrase
	body.TrackId = c.GetRespHeader(fiber.HeaderXRequestID)
	if err := c.BodyParser(&body); err != nil {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"error parsing request body",
			slog.Any("error", err.Error()),
		)
		return fiber.ErrUnprocessableEntity
	}

	if errs, ok := common.GetValidator().ValidateStruct(&body); !ok {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"validation error",
			slog.String("error", errs[0].Message),
		)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	result := usc.urlShortenerService.ShortUrlByCustomPhrase(&body)
	return c.Status(result.StatusCode).JSON(result)
}
