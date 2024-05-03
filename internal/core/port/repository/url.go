package repository

import (
	"errors"

	"github.com/artm2000/urlbook/internal/core/dto"
)

var (
	DUPLICATE_URL_PHRASE = errors.New("duplicated url phrase")
)

type Url interface {
	Insert(newUrl *dto.URL) error
}
