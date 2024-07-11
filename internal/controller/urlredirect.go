package controller

import (
	"context"
	"log/slog"
	"time"

	"github.com/artm2000/urlbook/internal/core/common"
	"github.com/artm2000/urlbook/internal/core/dto"
	"github.com/artm2000/urlbook/internal/core/model/request"
	"github.com/artm2000/urlbook/internal/core/port/service"
	"github.com/gofiber/fiber/v2"
)

type urlRedirect struct {
	urlShortenerService service.UrlShortener
	urlMetricsService   service.UrlMetrics
}

func NewUrlRedirect(urlShortenerService service.UrlShortener, urlMetricsService service.UrlMetrics) HttpController {
	return &urlRedirect{
		urlShortenerService,
		urlMetricsService,
	}
}

func (ur *urlRedirect) InitRoutes(api fiber.Router) {
	api.Get("/:urlRedirect", ur.redirectUrl)
}

func (ur *urlRedirect) GetPrefix() string { return "/" }

func (ur *urlRedirect) redirectUrl(c *fiber.Ctx) error {
	var params request.RedirectToDestination
	params.TrackId = c.GetRespHeader(fiber.HeaderXRequestID)
	if err := c.ParamsParser(&params); err != nil {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"error parsing request body",
			slog.Any("error", err.Error()),
		)
		return fiber.ErrUnprocessableEntity
	}

	if errs, ok := common.GetValidator().ValidateStruct(&params); !ok {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"validation error",
			slog.String("error", errs[0].Message),
		)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	result := ur.urlShortenerService.GetDestinationFromShortPhrase(&params)
	slog.Debug("redirecting to ...", "short_phrase", params.ShortPhrase, "destination", result.Data.DestinationUrl)

	redirectMetric := &dto.RedirectMetrics{
		ShortPhrase: params.ShortPhrase,
		UserAgent: c.GetReqHeaders()[fiber.HeaderUserAgent][0],
		IP: c.IP(),
		Time: time.Now(),
	}

	go ur.urlMetricsService.SubmitEvent(redirectMetric)

	return c.Status(fiber.StatusTemporaryRedirect).Redirect(result.Data.DestinationUrl)
}
