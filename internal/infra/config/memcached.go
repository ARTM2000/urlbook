package config

import (
	"github.com/bradfitz/gomemcache/memcache"
)

func NewMemcachedClient(addrs ...string) *memcache.Client {
	return memcache.New(addrs...)
}
