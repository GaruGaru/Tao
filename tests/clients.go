package tests

import (
	"github.com/go-redis/redis"
	"testing"
)

func TestRedisClient(t *testing.T) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: EnvOrDefault("REDIS_HOST", "localhost:6379"),
	})
}
