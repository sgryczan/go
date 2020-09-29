package redis

import (
	redis "github.com/go-redis/redis"
)

// Backend provides access to Redis
type Backend struct {
	client *redis.Client
}
