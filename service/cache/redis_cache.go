package cache

import (
	"fmt"
	"net"
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/service/props"
)

type RedisCache struct {
	address string
}

func (r *RedisCache) Get(key string) (map[string][]interface{}, bool) {

	conn, err := r.connect()
	if err != nil {
		return nil, false
	}
	defer conn.Close()

	return map[string][]interface{}{}, true
}

func (r *RedisCache) Set(key string, value map[string][]interface{}, ttl time.Duration) {

	conn, err := r.connect()
	if err != nil {
		return
	}
	defer conn.Close()
}

func (r *RedisCache) connect() (net.Conn, error) {
	return net.Dial("tcp", r.address)
}

func NewRedisCache(cfg props.CacheConfig) *RedisCache {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	return &RedisCache{address: addr}
}
