package service

import (
	"github.com/artm2000/urlbook/internal/core/dto"
	"github.com/artm2000/urlbook/internal/core/port/repository"
	"github.com/artm2000/urlbook/internal/core/port/service"
)

type urlMetrics struct {
	urlMetricsRepository repository.UrlMetrics
}

func NewUrlMetrics(urlMetricsRepository repository.UrlMetrics) service.UrlMetrics {
	return &urlMetrics{
		urlMetricsRepository,
	}
}

func (um *urlMetrics) SubmitEvent(event *dto.RedirectMetrics) error {
	err := um.urlMetricsRepository.InsertEvent(event)
	return err
}
