package supervisor

import (
	"fmt"
	"github.com/GaruGaru/Tao/tests"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/go-redis/redis"
	"math/rand"
	"testing"
)

func TestFileSystemLockObtainFail(t *testing.T) {

	key := fmt.Sprintf("test_storage_cleaner_%d", rand.Int31())
	redisClient := tests.TestRedisClient(t)

	redisClient.GeoAdd(key, &redis.GeoLocation{
		Name:      "expired",
		Latitude:  1.0,
		Longitude: 1.0,
	})

	cleaner := RedisStorageCleaner{
		Redis:     *redisClient,
		EventsKey: key,
		Statter:   &statsd.NoopClient{},
	}

	result, err := cleaner.Cleanup()

	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if result.Removed == 0 {
		t.Log("Expected to cleanup 1 element but got 0")
		t.FailNow()
	}

	res, err := redisClient.ZRange(cleaner.EventsKey, 0, -1).Result()

	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if len(res) != 0 {
		t.Log("Expected events zset to be empty")
		t.FailNow()
	}
}
