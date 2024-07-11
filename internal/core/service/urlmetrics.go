package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/artm2000/urlbook/internal/core/dto"
	"github.com/artm2000/urlbook/internal/core/port/repository"
	"github.com/artm2000/urlbook/internal/core/port/service"
)

type urlMetrics struct {
	urlMetricsRepository repository.UrlMetrics
	events               []*dto.RedirectMetrics
	mu                   sync.Mutex
}

func NewUrlMetrics(urlMetricsRepository repository.UrlMetrics) service.UrlMetrics {
	um := &urlMetrics{
		urlMetricsRepository: urlMetricsRepository,
		events: []*dto.RedirectMetrics{},
	}

	go func() {
		tk := time.NewTicker(time.Second * 40)
		for {
			<-tk.C
			if err := um.batchSubmitEvents(); err != nil {
				slog.LogAttrs(
					context.Background(), 
					slog.LevelError, 
					"fail to batch submit events", 
					slog.Any("error", err),
				)
			}
		}
	}()

	return um
}

func (um *urlMetrics) SubmitEvent(event *dto.RedirectMetrics) error {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.events = append(um.events, event)
	return nil
}

func (um *urlMetrics) batchSubmitEvents() error {
	um.mu.Lock()
	defer um.mu.Unlock()

	if len(um.events) == 0 {
		return nil
	}
	
	err := um.urlMetricsRepository.BatchInsertEvents(um.events)
	if err == nil {
		// in case that we don't have any error, clean the slice
		um.events = []*dto.RedirectMetrics{}
	}

	return err
}
