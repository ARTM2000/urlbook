package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/artm2000/urlbook/internal/core/dto"
	"github.com/artm2000/urlbook/internal/core/port/repository"
	"gorm.io/gorm"
)

type urlMetrics struct {
	ID          uint64    `gorm:"primaryKey;not null" json:"id"`
	ShortPhrase string    `gorm:"type:string;not null;" json:"short_phrase"`
	UserAgent   string    `gorm:"type:string;" json:"user_agent"`
	IP          string    `gorm:"type:string;not null;" json:"ip"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
}

type urlMetricsRepository struct {
	db *gorm.DB
}

func NewUrlMetricsRepository(db *gorm.DB) repository.UrlMetrics {
	db.AutoMigrate(&urlMetrics{})

	return &urlMetricsRepository{
		db,
	}
}

func (umr *urlMetricsRepository) InsertEvent(event *dto.RedirectMetrics) error {
	newUrlMetric := urlMetrics{
		ShortPhrase: event.ShortPhrase,
		UserAgent: event.UserAgent,
		IP: event.IP,
		CreatedAt: event.Time,
	}

	dbResult := umr.db.Model(&urlMetrics{}).Create(&newUrlMetric)
	if dbResult.Error != nil {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"unrecognized error on insert new url happened",
			slog.Any("error", dbResult.Error),
			slog.Any("event", event),
		)
		return repository.ErrNotRecognized
	}

	return nil
}
