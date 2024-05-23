package service

import (
	"time"

	"github.com/artm2000/urlbook/internal/core/dto"
	"github.com/artm2000/urlbook/internal/core/model/request"
	"github.com/artm2000/urlbook/internal/core/model/response"
)

type UrlShortener interface {
	ShortUrl(request *request.SubmitUrl) *response.Response[response.SubmitUrl]
	ShortUrlByCustomPhrase(request *request.SubmitUrlByCustomPhrase) *response.Response[response.SubmitUrl]
	GetDestinationFromShortPhrase(request *request.RedirectToDestination) *response.Response[response.GetUrlFromPhrase]
	GetManyShortUrls(from time.Time, to time.Time, size, page int) (*[]dto.URL, error)
}
