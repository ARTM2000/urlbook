package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/artm2000/urlbook/internal/core/dto"
	"github.com/artm2000/urlbook/internal/core/port/repository"
	"gorm.io/gorm"
)

type url struct {
	ID          uint64    `gorm:"primaryKey;not null" json:"id"`
	ShortPhrase string    `gorm:"type:string;not null;unique;uniqueIndex" json:"short_phrase"`
	Destination string    `gorm:"type:varchar(512);not null" json:"destination"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli;not null" json:"created_at"`
}

type urlRepository struct {
	db *gorm.DB
}

func NewUrlRepository(db *gorm.DB) repository.Url {
	// todo: consider better place
	db.AutoMigrate(&url{})

	return &urlRepository{
		db,
	}
}

func (ur *urlRepository) Insert(newUrl *dto.URL) error {
	newUrlModel := url{
		ShortPhrase: newUrl.ShortPhrase,
		Destination: newUrl.Destination,
	}
	dbResult := ur.db.Model(&url{}).Create(&newUrlModel)
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrDuplicatedKey) {
			slog.LogAttrs(
				context.Background(),
				slog.LevelDebug,
				fmt.Sprintf("no url found by short phrase %s", newUrlModel.ShortPhrase),
				slog.String("short_phrase", newUrlModel.ShortPhrase),
			)
			return repository.ErrDuplicateUrlPhrase
		}

		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"unrecognized error on insert new url happened",
			slog.Any("error", dbResult.Error),
			slog.String("short_phrase", newUrlModel.ShortPhrase),
		)
		return repository.ErrNotRecognized
	}

	return nil
}

func (ur *urlRepository) FindUrlByShortPhrase(shortPhrase string) (*dto.URL, error) {
	var foundUrl url
	dbResult := ur.db.Model(&url{}).Where(url{ShortPhrase: shortPhrase}).First(&foundUrl)
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			slog.LogAttrs(
				context.Background(),
				slog.LevelDebug,
				fmt.Sprintf("no url found by short phrase %s", shortPhrase),
				slog.String("short_phrase", shortPhrase),
			)
			return nil, repository.ErrNotFoundUrlPhrase
		}

		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"unrecognized error on finding url by short phrase happened",
			slog.Any("error", dbResult.Error),
			slog.String("short_phrase", shortPhrase),
		)
		return nil, repository.ErrNotRecognized
	}

	return &dto.URL{
		ShortPhrase: foundUrl.ShortPhrase,
		Destination: foundUrl.Destination,
		CreatedAt:   foundUrl.CreatedAt,
	}, nil
}
