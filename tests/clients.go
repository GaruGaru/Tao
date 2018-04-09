package tests

import (
	"testing"
	"github.com/go-redis/redis"
)

func TestRedisClient(t *testing.T) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: EnvOrDefault("REDIS_HOST","localhost:6379"),
	})
}


