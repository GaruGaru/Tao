package providers

import (
	"strconv"
	"github.com/cactus/go-statsd-client/statsd"
	"fmt"
	"time"
)

type CachedEventsProvider struct {
	Provider EventProvider
	statsd   statsd.Statter
	Cache    []EventsCache
}

func NewCachedEventsProvider(provider EventProvider, caches []EventsCache, statsd statsd.Statter) CachedEventsProvider {
	return CachedEventsProvider{
		Provider: provider,
		statsd:   statsd,
		Cache:    caches,
	}
}

func (p CachedEventsProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {

	key := requestKey(lat, lon, rng, sorting)

	for _, cache := range p.Cache {
		fmt.Printf("Fetching cache from: " + cache.Name())
		cacheFetchStart := time.Now()
		result, present, err := cache.Get(lat, lon, rng, sorting)
		if err == nil && present {
			fmt.Printf("Cache hit from: " + cache.Name())
			p.statsd.Inc(fmt.Sprintf("cache.%s.hit", cache.Name()), 1, 1.0)
			p.statsd.TimingDuration(fmt.Sprintf("cache.%s.latency", cache.Name()), time.Now().Sub(cacheFetchStart), 1.0)
			return result, nil
		} else if err != nil {
			p.statsd.Inc(fmt.Sprintf("cache.%s.error", cache.Name()), 1, 1.0)
		}
		p.statsd.Inc(fmt.Sprintf("cache.%s.miss", cache.Name()), 1, 1.0)
	}

	fmt.Printf("Cache miss")

	fetchStart := time.Now()

	events, err := p.Provider.Events(lat, lon, rng, sorting)

	p.statsd.TimingDuration("fetch.latency", time.Now().Sub(fetchStart), 1.0)

	if err != nil {
		p.statsd.Inc("fetch.error", 1, 1.0)
		return nil, err
	}

	p.statsd.Inc("fetch.ok", 1, 1.0)

	p.updateCache(events, key)

	return events, nil

}

func (p CachedEventsProvider) updateCache(events []DojoEvent, key string) {
	for _, cache := range p.Cache {
		cacheUpdateStart := time.Now()
		err := cache.Put(key, events)
		if err != nil {
			println(err.Error())
			p.statsd.Inc(fmt.Sprintf("cache.%s.update.error", cache.Name()), 1, 1.0)
		} else {
			p.statsd.Inc(fmt.Sprintf("cache.%s.update.ok", cache.Name()), 1, 1.0)
			p.statsd.TimingDuration(fmt.Sprintf("cache.%s.update.latency", cache.Name()), time.Now().Sub(cacheUpdateStart), 1.0)
		}
	}
}

func float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'E', -1, 64)
}
