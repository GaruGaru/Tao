package providers

import (
	"testing"
	"time"
	"github.com/cactus/go-statsd-client/statsd"

)

type TestEventProvider struct {
	events []DojoEvent
}

func (e TestEventProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {
	return e.events, nil
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

	if len(events) != len(testEvents){
		t.Log("Cached event count <> original events count")
		t.FailNow()
	}


}
