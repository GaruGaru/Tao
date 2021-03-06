package providers

import (
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"log"
	"strconv"
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

func (p CachedEventsProvider) fetchCache(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, bool) {
	key := RequestKey(lat, lon, rng, sorting)

	var cached []DojoEvent
	hitCacheIndex := 0
	hit := false

	for i, cache := range p.Cache {

		log.Println(fmt.Sprintf("Fetching cache from: %s", cache.Name()))
		cacheFetchStart := time.Now()
		result, present, err := cache.Get(lat, lon, rng, sorting)
		if err == nil && present {
			log.Println(fmt.Sprintf("Cache hit from: %s", cache.Name()))
			p.statsd.Inc(fmt.Sprintf("cache.%s.hit", cache.Name()), 1, 1.0)
			p.statsd.TimingDuration(fmt.Sprintf("cache.%s.latency", cache.Name()), time.Now().Sub(cacheFetchStart), 1.0)
			cached = result
			hitCacheIndex = i
			hit = true
			break
		} else if err != nil {
			log.Println(err.Error())
			p.statsd.Inc(fmt.Sprintf("cache.%s.error", cache.Name()), 1, 1.0)
		} else {
			p.statsd.Inc(fmt.Sprintf("cache.%s.miss", cache.Name()), 1, 1.0)
			log.Println(fmt.Sprintf("Cache miss from: %s", cache.Name()))
		}
	}

	if hit {
		for i := hitCacheIndex - 1; i >= 0; i-- {
			go p.updateCache(p.Cache[i], key, cached)
		}
	}

	return cached, hit
}

func (p CachedEventsProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {

	key := RequestKey(lat, lon, rng, sorting)

	cached, present := p.fetchCache(lat, lon, rng, sorting)

	if present {
		return cached, nil
	}

	fetchStart := time.Now()

	events, err := p.Provider.Events(lat, lon, rng, sorting)

	p.statsd.TimingDuration("fetch.latency", time.Now().Sub(fetchStart), 1.0)

	if err != nil {
		p.statsd.Inc("fetch.error", 1, 1.0)
		return nil, err
	}

	p.statsd.Inc("fetch.ok", 1, 1.0)

	for _, cache := range p.Cache {
		go p.updateCache(cache, key, events)
	}

	return events, nil

}

func (p CachedEventsProvider) updateCache(cache EventsCache, key string, events []DojoEvent) {
	cacheUpdateStart := time.Now()
	err := cache.Put(key, events)
	if err != nil {
		println(err.Error())
		p.statsd.Inc(fmt.Sprintf("cache.%s.update.error", cache.Name()), 1, 1.0)
	} else {
		log.Println("Updated cache " + cache.Name())
		p.statsd.Inc(fmt.Sprintf("cache.%s.update.ok", cache.Name()), 1, 1.0)
		p.statsd.TimingDuration(fmt.Sprintf("cache.%s.update.latency", cache.Name()), time.Now().Sub(cacheUpdateStart), 1.0)
	}
}

func float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'E', -1, 64)
}
