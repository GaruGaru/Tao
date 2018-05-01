package supervisor

import (
	"github.com/go-redis/redis"
	"github.com/cactus/go-statsd-client/statsd"
)

type CleanerResult struct {
	Removed int
}

type StorageCleaner interface {
	Cleanup() (CleanerResult, error)
}

type RedisStorageCleaner struct {
	Redis     redis.Client
	EventsKey string
	Statter   statsd.Statter
}

type NoOpCleaner struct {
}

func (c NoOpCleaner) Cleanup() (CleanerResult, error) {
	return CleanerResult{}, nil
}

func (cleaner RedisStorageCleaner) Cleanup() (CleanerResult, error) {
	cmd := cleaner.Redis.ZRange(cleaner.EventsKey, 0, -1)
	if cmd.Err() != nil {
		return CleanerResult{}, cmd.Err()
	}

	result, err := cmd.Result()

	if err != nil {
		cleaner.Statter.Inc("cleaner.run.error", 1, 1.0)
		return CleanerResult{}, err
	}

	removed := 0

	for _, k := range result {
		ecmd := cleaner.Redis.Exists(k)

		if ecmd.Err() != nil {
			return CleanerResult{}, ecmd.Err()
		}

		res, err := ecmd.Result()

		if err != nil {
			return CleanerResult{}, err
		}

		if res == 0 {
			removed++
			cleaner.Redis.ZRem(cleaner.EventsKey, k)
			cleaner.Statter.Inc("cleaner.result.removed", 1, 1.0)
		}

	}

	cleaner.Redis.BgSave()

	cleaner.Statter.Inc("cleaner.run.ok", 1, 1.0)

	return CleanerResult{Removed: removed,}, nil

}
