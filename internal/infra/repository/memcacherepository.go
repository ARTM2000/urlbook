package repository

import (
	"errors"
	"log/slog"
	"time"

	"github.com/artm2000/urlbook/internal/core/port/repository"
	"github.com/bradfitz/gomemcache/memcache"
)

type memcachedRepository struct {
	memcachedClient *memcache.Client
}

func NewMemcachedRepository(mc *memcache.Client) repository.Cache {
	return &memcachedRepository{
		memcachedClient: mc,
	}
}

func (mr *memcachedRepository) Get(key string) ([]byte, error) {
	item, err := mr.memcachedClient.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil, repository.ErrCacheMissed
		}

		return nil, repository.ErrCacheServerFault
	}

	return item.Value, nil
}

func (mr *memcachedRepository) Set(key string, value []byte, ttl time.Duration) error {
	item := memcache.Item{
		Key:   key,
		Value: value,
		Expiration: int32(ttl.Seconds()),
	}
	if err := mr.memcachedClient.Set(&item); err != nil {
		slog.Debug(err.Error())
		return repository.ErrCacheServerFault
	}

	return nil
}

func (mr *memcachedRepository) Del(key string) error {
	if err := mr.memcachedClient.Delete(key); err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return repository.ErrCacheMissed
		}

		return repository.ErrCacheServerFault
	}

	return nil
}
