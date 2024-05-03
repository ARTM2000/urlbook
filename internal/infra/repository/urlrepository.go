package repository

import (
	"time"

	"github.com/artm2000/urlbook/internal/core/dto"
	"github.com/artm2000/urlbook/internal/core/port/repository"
	"gorm.io/gorm"
)

type urlSchema struct {
	ID          uint64    `gorm:"primaryKey;not null" json:"id"`
	ShortPhrase string    `gorm:"type:string;not null;unique" json:"short_phrase"`
	Destination string    `gorm:"type:text;not null;unique" json:"destination"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli;not null" json:"created_at"`
}

type urlRepository struct {
	db *gorm.DB
}

func NewUrlRepository(db *gorm.DB) repository.Url {
	// todo: consider better place
	// db.AutoMigrate(&urlSchema{})

	return &urlRepository{
		db,
	}
}

func (ur *urlRepository) Insert(newUrl *dto.URL) error {
	// todo: develope the insertion process

	return nil
}
