package dto

import "time"

type RedirectMetrics struct {
	ShortPhrase string
	UserAgent   string
	IP          string
	Time        time.Time
}
