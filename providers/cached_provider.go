package providers

import (
	"strings"
	"strconv"
	"time"
	statsd2 "github.com/smira/go-statsd"
)

type CachedEventsProvider struct {
	Provider EventProvider
	cache    map[string][]DojoEvent
	statsd    statsd2.Client
}

func NewCachedEventsProvider(provider EventProvider, statsd statsd2.Client) CachedEventsProvider {
	return CachedEventsProvider{
		Provider: provider,
		cache:    make(map[string][]DojoEvent),
		statsd:   statsd,
	}
}

func (p CachedEventsProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {

	key := requestKey(lat, lon, rng, sorting)

	cached, cachePresent := p.cache[key]

	if cachePresent {
		p.statsd.Incr("cache.hit",1)
		return cached, nil
	} else {
		p.statsd.Incr("cache.miss",1)
		t1 := time.Now()

		events, err := p.Provider.Events(lat, lon, rng, sorting)

		t2 := time.Now()
		duration := int64(t2.Sub(t1) / time.Millisecond)
		p.statsd.Timing("api.latency", duration)

		if err == nil {
			p.statsd.Incr("api.ok",1)
			p.cache[key] = events
			return events, nil
		} else {
			p.statsd.Incr("api.fail",1)
			return nil, err
		}

	}

}

func requestKey(lat float64, lon float64, rng int, sorting string) string {
	return strings.Join(
		[]string{float64ToString(lat), float64ToString(lon), strconv.Itoa(rng), sorting},
		":",
	)
}

func float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'E', -1, 64)
}
