package service

import "github.com/artm2000/urlbook/internal/core/dto"

type UrlMetrics interface {
	SubmitEvent(metric *dto.RedirectMetrics) error
}
