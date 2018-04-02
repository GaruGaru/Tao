package providers

import (
	"github.com/cactus/go-statsd-client/statsd"
	"testing"
	"time"
	"fmt"
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
