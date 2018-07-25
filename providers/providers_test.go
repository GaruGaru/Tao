package providers

import (
	"encoding/json"
	"fmt"
	"github.com/GaruGaru/Tao/tests"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/go-redis/redis"
	"math/rand"
	"testing"
	"time"
)

type TestEventProvider struct {
	events []DojoEvent
}

func (e TestEventProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {
	return e.events, nil
}

type FailingEventProvider struct {
	events []DojoEvent
}

func (e FailingEventProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {
	return []DojoEvent{}, fmt.Errorf("generic error")
}

func TestAllFailingProvider(t *testing.T) {

	provider := ProvidersManager{Providers: []EventProvider{
		FailingEventProvider{},
		FailingEventProvider{},
	}}

	events, err := provider.Events(1.0, 1.0, 1, "distance")

	if err != nil {
		t.Log("Error fetching from cache")
		t.FailNow()
	}

	if len(events) != 0 {
		t.Log("Cached event count <> original events count")
		t.FailNow()
	}

}

func eventToJson(e DojoEvent) string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestRedisEventsProvider(t *testing.T) {

	prefix := rand.Int()
	geoKey := fmt.Sprintf("%d_test_locations", prefix)
	redisClient := tests.TestRedisClient(t)
	testEvents := loadTestData(t)

	for _, e := range testEvents {
		key := fmt.Sprintf("%d:%s:%s:%d", prefix, e.Title, e.TicketUrl, e.StartTime)
		geoResult := redisClient.GeoAdd(geoKey, &redis.GeoLocation{Longitude: e.Location.Longitude, Latitude: e.Location.Latitude, Name: key})
		redisClient.Set(key, eventToJson(e), 0)
		if geoResult.Err() != nil {
			t.FailNow()
		}
	}

	redisProvider := RedisEventsProvider{
		Redis:        *redisClient,
		LocationsKey: geoKey,
	}

	events, err := redisProvider.Events(0, 0, 1000000, "distance")

	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if len(events) != len(testEvents) {
		t.Log(fmt.Sprintf("Expecting %d events but got %d", len(testEvents), len(events)))
		t.FailNow()
	}

}

func TestFailingProvider(t *testing.T) {

	testEvents := loadTestData(t)

	provider := ProvidersManager{Providers: []EventProvider{
		TestEventProvider{events: testEvents},
		TestEventProvider{events: testEvents},
		FailingEventProvider{},
		FailingEventProvider{},
	}}

	events, err := provider.Events(1.0, 1.0, 1, "distance")

	if err != nil {
		t.Log("Error fetching from cache")
		t.FailNow()
	}

	if len(events) != len(testEvents)*2 {
		t.Log("Cached event count <> original events count")
		t.FailNow()
	}

}

func TestCachedProvider(t *testing.T) {

	testEvents := loadTestData(t)

	provider := ProvidersManager{Providers: []EventProvider{TestEventProvider{events: testEvents}}}

	caches := []EventsCache{
		NewLocalCache(1 * time.Minute),
	}

	cachedProvider := NewCachedEventsProvider(provider, caches, &statsd.NoopClient{})

	events, err := cachedProvider.Events(1.0, 1.0, 1, "distance")

	if err != nil {
		t.Log("Error fetching from cache")
		t.FailNow()
	}

	if len(events) != len(testEvents) {
		t.Log("Cached event count <> original events count")
		t.FailNow()
	}

}
