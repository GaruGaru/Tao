package providers

import (
	"strconv"
	"github.com/cactus/go-statsd-client/statsd"
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
		result, present, err := cache.Get(lat, lon, rng, sorting)
		if err == nil && present {
			p.statsd.Inc("cache.hit."+cache.Name(), 1, 1.0)
			return result, nil
		} else if err != nil {
			p.statsd.Inc("cache.error."+cache.Name(), 1, 1.0)
		}
		p.statsd.Inc("cache.miss." + cache.Name(), 1, 1.0)
	}

	events, err := p.Provider.Events(lat, lon, rng, sorting)

	if err != nil {
		return nil, err
	}

	p.updateCache(events, key)

	return events, nil

}

func (p CachedEventsProvider) updateCache(events []DojoEvent, key string) {
	for _, cache := range p.Cache {
		err := cache.Put(key, events)
		if err != nil {
			println(err.Error())
			p.statsd.Inc("cache.update_error." + cache.Name(), 1, 1.0)
		}else{
			p.statsd.Inc("cache.update." + cache.Name(), 1, 1.0)
		}
	}
}

func float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'E', -1, 64)
}
