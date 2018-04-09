package cmd

import (
	"github.com/spf13/viper"
	"github.com/go-redis/redis"
	"github.com/GaruGaru/Tao/scraper"
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"time"
)

func GetConfiguredStorage() string {
	return viper.GetString("storage")
}

func GetStatter() statsd.Statter {
	host := viper.GetString("statsd_host")

	if host == "" {
		return &statsd.NoopClient{}
	}

	client, err := statsd.NewClient(host, "tao")

	if err != nil {
		return &statsd.NoopClient{}
	}

	return client
}

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis_host"),
	})
}

func GetRedisLocationsKey() string {
	return "tao_events_locations"
}

func GetScraperLock() scraper.DojoScraperLock {
	storage := GetConfiguredStorage()
	if storage == "local" {
		return scraper.FileSystemLock{LockFile: "/tmp/tao.lock"}
	} else if storage == "redis" {
		return scraper.RedisDojoScraperLock{Redis: *GetRedisClient(), LockKey: "tao_lock"}
	}
	panic(fmt.Errorf("unable to create scraper lock instance for storage type: %s", storage))
}

func GetScraperStorage() scraper.EventsStorage {
	storage := GetConfiguredStorage()
	if storage == "local" {
		return scraper.FileSystemEventsStorage{StoreFile: "events.json"}
	} else if storage == "redis" {
		return scraper.RedisEventsStorage{
			Redis:  *GetRedisClient(),
			GeoKey: GetRedisLocationsKey(),
		}
	}
	panic(fmt.Errorf("unable to create scraper events storage instance for storage type: %s", storage))
}

func GetScraperDelayer() scraper.ScraperDelayer {
	storage := GetConfiguredStorage()
	if storage == "local" {
		return scraper.LocalScraperDelayer{}
	} else if storage == "redis" {
		return scraper.NewRedisScraperDelayer(*GetRedisClient(), time.Duration(viper.GetInt("scraper_delay"))*time.Second)
	}
	panic(fmt.Errorf("unable to create scraper events delayer instance for storage type: %s", storage))
}
