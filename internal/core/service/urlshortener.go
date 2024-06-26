package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

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
		err := ush.cacheRepo.Set(fmt.Sprintf("url:%s", shortPhrase), []byte(request.Url), repository.NewUrlCacheTTL)
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

func (ush *urlShortener) ShortUrlByCustomPhrase(request *request.SubmitUrlByCustomPhrase) *response.Response[response.SubmitUrl] {
	ctx := context.Background()

	if _, err := ush.urlRepo.FindUrlByShortPhrase(request.Phrase); err == nil {
		// in this case, there is a same short-phrase existing for a url
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"there is a existing short-phrase for this custom short-phrase",
			slog.String("short_phrase", request.Phrase),
		)
		return createFailResponse[response.SubmitUrl](message.DUPLICATE_SHORT_PHRASE, request.TrackId, code.CONFLICT)
	}

	urlDto := dto.URL{
		ShortPhrase: request.Phrase,
		Destination: request.Url,
	}
	if err := ush.urlRepo.Insert(&urlDto); err != nil {
		if errors.Is(err, repository.ErrDuplicateUrlPhrase) {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"request failed with DUPLICATE_URL_PHRASE",
				slog.String("destination", request.Url),
				slog.String("short_phrase", request.Phrase),
			)
			return createFailResponse[response.SubmitUrl](message.INTERNAL_SYSTEM_ERROR, request.TrackId, code.INTERNAL_SYSTEM_ERROR)
		}

		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"request failed",
			slog.String("destination", request.Url),
			slog.String("short_phrase", request.Phrase),
			slog.Any("error", err.Error()),
		)
		return createFailResponse[response.SubmitUrl](message.INTERNAL_SYSTEM_ERROR, request.TrackId, code.INTERNAL_SYSTEM_ERROR)
	}

	go func() {
		err := ush.cacheRepo.Set(fmt.Sprintf("url:%s", request.Phrase), []byte(request.Url), repository.NewUrlCacheTTL)
		slog.LogAttrs(
			ctx,
			slog.LevelDebug,
			"set short url cache",
			slog.String("short_phrase", request.Phrase),
			slog.String("destination", request.Url),
			slog.Any("error", err),
		)
	}()

	submitUrlRes := response.SubmitUrl{
		ShortUrl: ush.createFullShortUrl(request.Phrase),
	}
	return createSuccessResponse(&submitUrlRes, request.TrackId)
}

func (ush *urlShortener) GetDestinationFromShortPhrase(request *request.RedirectToDestination) *response.Response[response.GetUrlFromPhrase] {
	ctx := context.Background()

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
				ctx,
				slog.LevelError,
				"request failed with ErrNotFoundUrlPhrase",
				slog.String("short_phrase", request.ShortPhrase),
			)
			return createFailResponse[response.GetUrlFromPhrase](message.NOT_FOUND_URL, request.TrackId, code.NOT_FOUND)
		}

		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"request failed",
			slog.String("short_phrase", request.ShortPhrase),
			slog.Any("error", err.Error()),
		)
		return createFailResponse[response.GetUrlFromPhrase](message.INTERNAL_SYSTEM_ERROR, request.TrackId, code.INTERNAL_SYSTEM_ERROR)
	}

	go func() {
		err := ush.cacheRepo.Set(fmt.Sprintf("url:%s", request.ShortPhrase), []byte(url.Destination), repository.UrlDefaultCacheTTL)
		slog.LogAttrs(
			ctx,
			slog.LevelDebug,
			"set short url cache",
			slog.String("short_phrase", request.ShortPhrase),
			slog.String("destination", url.Destination),
			slog.Any("error", err),
		)
	}()

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

func (ush *urlShortener) GetManyShortUrls(from time.Time, to time.Time, size, page int) (*[]dto.URL, error) {
	urls, err := ush.urlRepo.FindManyUrlsByTimeScope(from, to, size, page)
	if err != nil {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"error getting many short urls",
			slog.Time("from", from),
			slog.Time("to", to),
			slog.Int("size", size),
			slog.Int("page", page),
			slog.Any("error", err),
		)
		return nil, errors.New(string(message.INTERNAL_SYSTEM_ERROR))
	}
	return urls, nil
}
