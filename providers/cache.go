package providers

import (
	"github.com/go-redis/redis"
	"github.com/wunderlist/ttlcache"
	"encoding/json"
	"time"
	"strings"
	"strconv"
)

func requestKey(lat float64, lon float64, rng int, sorting string) string {
	return strings.Join(
		[]string{float64ToString(lat), float64ToString(lon), strconv.Itoa(rng), sorting},
		":",
	)
}

func eventsToJson(events []DojoEvent) (string, error) {
	eventsJson, err := json.Marshal(events)
	if err != nil {
		return "", err
	}
	return string(eventsJson), nil
}

func eventsFromJson(jsonStr string) ([]DojoEvent, error) {
	var events []DojoEvent
	err := json.Unmarshal([]byte(jsonStr), &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

type EventsCache interface {
	Get(lat float64, lon float64, rng int, sorting string) (events []DojoEvent, present bool, err error)
	Put(key string, events []DojoEvent) error
	Name() string
}

type RedisEventCache struct {
	redis         redis.Client
	cacheDuration time.Duration
}

func (rc RedisEventCache) Name() string {
	return "redis"
}

func (rc RedisEventCache) Get(lat float64, lon float64, rng int, sorting string) (events []DojoEvent, present bool, err error) {
	key := requestKey(lat, lon, rng, sorting)
	cache, err := rc.redis.Get(key).Result()
	if err != nil {
		return nil, false, err
	} else if err == redis.Nil {
		return nil, false, nil
	} else {
		events, err := eventsFromJson(cache)
		if err != nil {
			return nil, false, err
		}
		return events, true, nil
	}
}

func (rc RedisEventCache) Put(key string, events []DojoEvent) error {
	eventsJson, err := eventsToJson(events)
	if err != nil {
		return err
	}
	return rc.redis.Set(key, string(eventsJson), rc.cacheDuration).Err()
}

type LocalEventCache struct {
	cacheDuration time.Duration
	cache         ttlcache.Cache
}

func NewLocalCache(duration time.Duration) LocalEventCache {
	return LocalEventCache{
		cacheDuration: duration,
		cache:         *ttlcache.NewCache(duration),
	}
}

func (rc LocalEventCache) Name() string {
	return "memory"
}

func (rc LocalEventCache) Get(lat float64, lon float64, rng int, sorting string) (events []DojoEvent, present bool, err error) {
	key := requestKey(lat, lon, rng, sorting)
	cache, exists := rc.cache.Get(key)
	if exists {
		events, err := eventsFromJson(cache)
		if err != nil {
			return nil, false, err
		}
		return events, true, nil
	} else {
		return nil, false, nil
	}

}

func (rc LocalEventCache) Put(key string, events []DojoEvent) error {
	eventsJson, err := json.Marshal(events)
	if err != nil {
		return err
	}
	rc.cache.Set(key, string(eventsJson))
	return nil
}
