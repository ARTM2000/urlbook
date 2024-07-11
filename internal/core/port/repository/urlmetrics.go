package repository

import "github.com/artm2000/urlbook/internal/core/dto"

type UrlMetrics interface {
	InsertEvent(event *dto.RedirectMetrics) error
	BatchInsertEvents(events []*dto.RedirectMetrics) error
}
