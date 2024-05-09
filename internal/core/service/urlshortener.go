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
	cacheRepo    repository.Cache
	baseShortUrl string
}

func NewUrlShortener(urlRepo repository.Url, cacheRepo repository.Cache, baseShortUrl string) service.UrlShortener {
	return &urlShortener{
		urlRepo,
		cacheRepo,
		baseShortUrl,
	}
}

func (ush *urlShortener) ShortUrl(request *request.SubmitUrl) *response.Response[response.SubmitUrl] {
	// todo: add some validation on url
	ctx := context.Background()
	shortPhrase := ush.generateShortUrlPhrase(request.Url)
	urlDto := dto.URL{
		ShortPhrase: shortPhrase,
		Destination: request.Url,
	}

	err := ush.urlRepo.Insert(&urlDto)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateUrlPhrase) {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"request failed with DUPLICATE_URL_PHRASE",
				slog.String("destination", request.Url),
				slog.String("short_phrase", shortPhrase),
			)
			return createFailResponse[response.SubmitUrl](message.INTERNAL_SYSTEM_ERROR, request.TrackId, code.INTERNAL_SYSTEM_ERROR)
		}

		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"request failed",
			slog.String("destination", request.Url),
			slog.String("short_phrase", shortPhrase),
			slog.Any("error", err.Error()),
		)
		return createFailResponse[response.SubmitUrl](message.INTERNAL_SYSTEM_ERROR, request.TrackId, code.INTERNAL_SYSTEM_ERROR)
	}
	
	go func() {
		err := ush.cacheRepo.Set(fmt.Sprintf("url:%s", shortPhrase), []byte(request.Url))
		slog.LogAttrs(
			ctx, 
			slog.LevelDebug, 
			"set short url cache", 
			slog.String("short_phrase", shortPhrase),
			slog.String("destination", request.Url),
			slog.Any("error", err),
		)
	}()

	submitUrlRes := response.SubmitUrl{
		ShortUrl: ush.createFullShortUrl(shortPhrase),
	}
	return createSuccessResponse(&submitUrlRes, request.TrackId)
}

func (ush *urlShortener) GetDestinationFromShortPhrase(request *request.RedirectToDestination) *response.Response[response.GetUrlFromPhrase] {
	// First, try to get destination from cache
	if dest, _ := ush.cacheRepo.Get(fmt.Sprintf("url:%s", request.ShortPhrase)); dest != nil {
		slog.Debug("got destination from cache", "short_phrase", request.ShortPhrase, "destination", string(dest))
		getUrlFromPhrase := response.GetUrlFromPhrase{
			DestinationUrl: string(dest),
		}
		return createSuccessResponse(&getUrlFromPhrase, request.TrackId)
	}

	url, err := ush.urlRepo.FindUrlByShortPhrase(request.ShortPhrase)
	if err != nil {
		if errors.Is(err, repository.ErrNotFoundUrlPhrase) {
			slog.LogAttrs(
				context.Background(),
				slog.LevelError,
				"request failed with ErrNotFoundUrlPhrase",
				slog.String("short_phrase", request.ShortPhrase),
			)
			return createFailResponse[response.GetUrlFromPhrase](message.NOT_FOUND_URL, request.TrackId, code.NOT_FOUND)
		}

		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"request failed",
			slog.String("short_phrase", request.ShortPhrase),
			slog.Any("error", err.Error()),
		)
		return createFailResponse[response.GetUrlFromPhrase](message.INTERNAL_SYSTEM_ERROR, request.TrackId, code.INTERNAL_SYSTEM_ERROR)
	}

	getUrlFromPhrase := response.GetUrlFromPhrase{
		DestinationUrl: url.Destination,
	}
	return createSuccessResponse(&getUrlFromPhrase, request.TrackId)
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
