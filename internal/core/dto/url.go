package dto

import "time"

type URL struct {
	ShortPhrase string
	Destination string
	CreatedAt   time.Time
}
