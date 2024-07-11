package config

import (
	"log/slog"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
)

func NewMemcachedClient(addrs ...string) *memcache.Client {
	slog.Debug("memcached address", slog.String("address", strings.Join(addrs, "-")))
	return memcache.New(addrs...)
}
