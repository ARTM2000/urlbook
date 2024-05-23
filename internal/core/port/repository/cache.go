package repository

import (
	"errors"
	"time"
)

var (
	ErrCacheMissed      = errors.New("ErrCacheMissed")
	ErrCacheServerFault = errors.New("ErrCacheServerFault")
	ErrMalformedKey     = errors.New("ErrMalformedKey")
)

const (
	NewUrlCacheTTL     = time.Hour * 24 // a day
	UrlDefaultCacheTTL = time.Hour * 3
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
	Del(key string) error
}
