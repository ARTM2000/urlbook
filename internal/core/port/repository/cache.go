package repository

import "errors"

var (
	ErrCacheMissed      = errors.New("ErrCacheMissed")
	ErrCacheServerFault = errors.New("ErrCacheServerFault")
	ErrMalformedKey     = errors.New("ErrMalformedKey")
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Del(key string) error
}
