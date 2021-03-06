package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/GaruGaru/Tao/providers"
	"github.com/GaruGaru/Tao/tests"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/go-redis/redis"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func LoadTestEvents(t *testing.T, path string) []providers.DojoEvent {
	testJson, err := ioutil.ReadFile(path)

	if err != nil {
		t.Log("Unable to read test json file")
		t.FailNow()
	}

	var events []providers.DojoEvent

	err = json.Unmarshal(testJson, &events)

	if err != nil {
		t.Log("Unable to unmarshal test json file")
		t.FailNow()
	}

	return events
}

type TestEventProvider struct {
	DojoEvents []providers.DojoEvent
}

func (p TestEventProvider) Events(lat float64, lon float64, rng int, sorting string) ([]providers.DojoEvent, error) {
	return p.DojoEvents, nil
}

func TestEventsScraper(t *testing.T) {

	testEvents := LoadTestEvents(t, "testdata/provider_response.json")

	store := NewInMemoryEventsStorage()
	dojoScraper := DojoScraper{
		Storage: store,
		Scraper: DefaultEventScraper{Provider: providers.ProvidersManager{Providers: []providers.EventProvider{
			TestEventProvider{DojoEvents: testEvents},
			TestEventProvider{DojoEvents: testEvents},
			TestEventProvider{DojoEvents: testEvents},
			TestEventProvider{DojoEvents: testEvents},
			TestEventProvider{DojoEvents: testEvents},
		}}},
		Lock:    FileSystemLock{LockFile: fmt.Sprintf("test.lock.%d", rand.Int())},
		Delayer: LocalScraperDelayer{Delay: 0, lastRun: time.Now()},
		Statter: &statsd.NoopClient{},
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

	testEvents := LoadTestEvents(t, "testdata/provider_response.json")

	for i := range testEvents {
		e := &testEvents[i]
		e.StartTime = time.Now().Add(48 * time.Hour).Unix()
		e.EndTime = time.Now().Add(49 * time.Hour).Unix()
	}

	delayerKey := fmt.Sprintf("test_delayer_%d", rand.Int31())

	geoKey := fmt.Sprintf("locations_test_%d", rand.Int31())

	lockKey := fmt.Sprintf("test_lock_%d", rand.Int31())

	dojoScraper := DojoScraper{
		Storage: RedisEventsStorage{Redis: *redisClient, GeoKey: geoKey},
		Scraper: DefaultEventScraper{Provider: TestEventProvider{DojoEvents: testEvents}},
		Lock:    RedisDojoScraperLock{Redis: *redisClient, LockKey: lockKey},
		Delayer: RedisScraperDelayer{Redis: *redisClient, Delay: 60 * time.Second, TimeKey: delayerKey},
		Statter: &statsd.NoopClient{},
	}

	defer dojoScraper.Lock.Release()

	dojoScraper.Run()

	obtain := dojoScraper.Lock.Obtain()

	if obtain != nil {
		t.Log(obtain)
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

	key, err := redisClient.Get(delayerKey).Result()

	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if key == "" {
		t.Log("empty delayer key")
		t.FailNow()
	}

}
