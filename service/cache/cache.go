package cache

import "time"

// Cache defines a generic caching interface.
// You can implement Redis using this.
type Cache interface {
	Get(key string) (map[string][]interface{}, bool)
	Set(key string, value map[string][]interface{}, ttl time.Duration)
}
