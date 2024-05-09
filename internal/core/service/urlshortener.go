package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/artm2000/urlbook/internal/core/common"
	"github.com/artm2000/urlbook/internal/core/dto"
	"github.com/artm2000/urlbook/internal/core/entity/code"
	"github.com/artm2000/urlbook/internal/core/entity/message"
	"github.com/artm2000/urlbook/internal/core/model/request"
	"github.com/artm2000/urlbook/internal/core/model/response"
	"github.com/artm2000/urlbook/internal/core/port/repository"
	"github.com/artm2000/urlbook/internal/core/port/service"
)

type urlShortener struct {
	urlRepo      repository.Url
	baseShortUrl string
}

func NewUrlShortener(urlRepo repository.Url, baseShortUrl string) service.UrlShortener {
	return &urlShortener{
		urlRepo,
		baseShortUrl,
	}
}

func (ush *urlShortener) ShortUrl(request *request.SubmitUrl) *response.Response[response.SubmitUrl] {
	// todo: add some validation on url

	shortPhrase := ush.generateShortUrlPhrase(request.Url)
	urlDto := dto.URL{
		ShortPhrase: shortPhrase,
		Destination: request.Url,
	}

	err := ush.urlRepo.Insert(&urlDto)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateUrlPhrase) {
			slog.LogAttrs(
				context.Background(),
				slog.LevelError,
				"request failed with DUPLICATE_URL_PHRASE",
				slog.String("destination", request.Url),
				slog.String("short_phrase", shortPhrase),
			)
			return createFailResponse[response.SubmitUrl](message.INTERNAL_SYSTEM_ERROR, request.TrackId, code.INTERNAL_SYSTEM_ERROR)
		}

		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"request failed",
			slog.String("destination", request.Url),
			slog.String("short_phrase", shortPhrase),
			slog.Any("error", err.Error()),
		)
		return createFailResponse[response.SubmitUrl](message.INTERNAL_SYSTEM_ERROR, request.TrackId, code.INTERNAL_SYSTEM_ERROR)
	}

	submitUrlRes := response.SubmitUrl{
		ShortUrl: ush.createFullShortUrl(shortPhrase),
	}
	return createSuccessResponse(&submitUrlRes, request.TrackId)
}

func (ush *urlShortener) generateShortUrlPhrase(destination string) string {
	shortPhrase := common.GetRandomUrlShortPhrase()
	sameUrl, _ := ush.urlRepo.FindUrlByShortPhrase(shortPhrase)
	for sameUrl != nil {
		shortPhrase = common.GetRandomUrlShortPhrase()
		sameUrl, _ = ush.urlRepo.FindUrlByShortPhrase(shortPhrase)
	}

	slog.Debug("generate short url phrase", "destination", destination, "short_phrase", shortPhrase)
	return shortPhrase
}

func (ush *urlShortener) createFullShortUrl(phrase string) string {
	return fmt.Sprintf("%s/%s", ush.baseShortUrl, phrase)
}
