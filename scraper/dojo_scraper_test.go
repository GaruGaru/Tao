package scraper

import (
	"testing"
	"fmt"
	"math/rand"
	"github.com/GaruGaru/Tao/tests"
	"github.com/go-redis/redis"
)

func TestEventsScraper(t *testing.T) {

	testEvents := tests.LoadTestEvents(t, "testdata/provider_response.json")

	store := NewInMemoryEventsStorage()
	dojoScraper := DojoScraper{
		Storage: store,
		Scraper: DefaultEventScraper{Provider: tests.TestEventProvider{DojoEvents: testEvents},},
		Lock:    FileSystemLock{LockFile: fmt.Sprintf("test.lock.%d", rand.Int())},
	}

	dojoScraper.Run()

	if dojoScraper.Lock.Release() == nil {
		t.Log("Scraper hasn't released the lock")
		t.FailNow()
	}

	eventsCount := len(store.Storage)

	if eventsCount != len(testEvents) {
		t.Log(fmt.Sprintf("Expected %d events in the store but where %d", len(testEvents), eventsCount))
		t.FailNow()
	}

}

func TestEventsScraperWithRedis(t *testing.T) {

	redisClient := tests.TestRedisClient(t)

	testEvents := tests.LoadTestEvents(t, "testdata/provider_response.json")

	geoKey := fmt.Sprintf("locations_test_%d", rand.Int())

	dojoScraper := DojoScraper{
		Storage: RedisEventsStorage{Redis: *redisClient, GeoKey: geoKey},
		Scraper: DefaultEventScraper{Provider: tests.TestEventProvider{DojoEvents: testEvents},},
		Lock:    RedisDojoScraperLock{Redis: *redisClient, LockKey: fmt.Sprintf("test_lock_%d", rand.Int())},
	}

	dojoScraper.Run()

	obtain := dojoScraper.Lock.Obtain()

	if obtain != nil {
		t.Log("Scraper hasn't released the lock")
		t.FailNow()
	}

	locations := redisClient.GeoRadius(geoKey, 0, 0, &redis.GeoRadiusQuery{
		Radius: 1000000,
		Unit:   "km",
	})

	result, err := locations.Result()

	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if len(result) != len(testEvents) {
		t.Log("Event count mismatch")
		t.FailNow()
	}

	for _, r := range result {
		key := r.Name
		result, err := redisClient.Get(key).Result()

		if err != nil {
			t.Log(err.Error())
			t.FailNow()
		}

		if result == "" {
			t.Log("Empty result for key" + key)
			t.FailNow()
		}
	}

}
