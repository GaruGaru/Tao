package main

import (
	"github.com/GaruGaru/Tao/providers"
	"github.com/GaruGaru/Tao/api"
	"github.com/cactus/go-statsd-client/statsd"
	"os"
	"github.com/go-redis/redis"
	"time"
)

func main() {

	statsdClient, err := statsd.NewClient(getEnv("STATSD_HOST", "localhost:8125"), "tao")

	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_HOST", "localhost:6379"),
		Password: "",
		DB:       0,
	})

	localCacheExpiration := 15 * time.Minute
	remoteCacheExpiration := 30 * time.Minute

	caches := []providers.EventsCache{
		providers.NewLocalCache(localCacheExpiration),
		providers.NewRedisEventsCache(*redisClient, remoteCacheExpiration),
	}

	provider := providers.NewCachedEventsProvider(providers.ProvidersManager{
		Providers: []providers.EventProvider{providers.NewEventBriteProvider(),},
	}, caches, statsdClient)

	taoApi := api.EventsApi{Provider: provider, Statsd: statsdClient}

	taoApi.Run()

}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
