package service

import (
	"github.com/artm2000/urlbook/internal/core/model/request"
	"github.com/artm2000/urlbook/internal/core/model/response"
)

type UrlShortener interface {
	ShortUrl(request *request.SubmitUrl) *response.Response[response.SubmitUrl]
	GetDestinationFromShortPhrase(request *request.RedirectToDestination) *response.Response[response.GetUrlFromPhrase]
}
