package cache

import "time"

// NoopCache is a cache implementation that stores nothing.
// Useful for tests.
type NoopCache struct{}

func NewNoopCache() Cache {
	return &NoopCache{}
}

func (c *NoopCache) Get(key string) (map[string][]interface{}, bool) {
	return nil, false
}

func (c *NoopCache) Set(key string, value map[string][]interface{}, ttlSeconds time.Duration) {
	// no-op
}
