package scraper

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type ScraperDelayer interface {
	CanRun() (bool, error)
	Refresh() error
}

type LocalScraperDelayer struct {
	Delay   time.Duration
	lastRun time.Time
}

func (d LocalScraperDelayer) CanRun() (bool, error) {
	return time.Now().Sub(d.lastRun) > d.Delay, nil
}

func (d LocalScraperDelayer) Refresh() error {
	d.lastRun = time.Now()
	return nil
}

type RedisScraperDelayer struct {
	Redis   redis.Client
	Delay   time.Duration
	TimeKey string
}

func NewRedisScraperDelayer(redis redis.Client, delay time.Duration) ScraperDelayer {
	return RedisScraperDelayer{
		Redis:   redis,
		Delay:   delay,
		TimeKey: "scraper_last_run",
	}
}

func (d RedisScraperDelayer) CanRun() (bool, error) {
	res, err := d.Redis.Get(d.TimeKey).Result()

	if err == redis.Nil {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	timestamp, err := strconv.Atoi(res)

	if err != nil {
		return false, err
	}

	delay := time.Now().Sub(time.Unix(int64(timestamp), 0))

	return time.Duration(delay) > d.Delay, nil
}
func (d RedisScraperDelayer) Refresh() error {
	return d.Redis.Set(d.TimeKey, time.Now().Unix(), 0).Err()
}
