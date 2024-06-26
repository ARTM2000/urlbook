package repository

import (
	"errors"
	"time"

	"github.com/artm2000/urlbook/internal/core/dto"
)

var (
	ErrDuplicateUrlPhrase = errors.New("duplicate url phrase")
	ErrNotFoundUrlPhrase  = errors.New("url phrase not found")
	ErrNotRecognized      = errors.New("not recognized error")
)

type Url interface {
	Insert(newUrl *dto.URL) error
	FindUrlByShortPhrase(shortPhrase string) (*dto.URL, error)
	FindManyUrlsByTimeScope(from time.Time, to time.Time, size, page int) (*[]dto.URL, error)
}
