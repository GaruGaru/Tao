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

	statsdClient, _ := statsd.NewClient(getEnv("STATSD_HOST", "localhost:8125"), "tao")

	//if err != nil {
	//	panic(err)
	//}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_HOST", "localhost:6379"),
		Password: "",
		DB:       0,
	})

	cacheDuration := 10 * time.Minute

	provider := providers.NewCachedEventsProvider(providers.ProvidersManager{
		Providers: []providers.EventProvider{
			providers.NewEventBriteProvider(),
		},
	}, *redisClient, cacheDuration, statsdClient)

	taoApi := api.EventsApi{Provider: provider}

	taoApi.Run()

}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
