package providers

import (
	"strings"
	"strconv"
	statsd2 "github.com/smira/go-statsd"
	"github.com/go-redis/redis"
	"encoding/json"
	"time"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type CachedEventsProvider struct {
	Provider      EventProvider
	statsd        statsd2.Client
	redis         redis.Client
	cacheDuration time.Duration
}

func NewCachedEventsProvider(provider EventProvider, redis redis.Client, cacheDuration time.Duration, statsd statsd2.Client) CachedEventsProvider {
	return CachedEventsProvider{
		Provider:      provider,
		statsd:        statsd,
		redis:         redis,
		cacheDuration: cacheDuration,
	}
}

func (p CachedEventsProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {

	key := requestKey(lat, lon, rng, sorting)

	cache, err := p.redis.Get(key).Result()

	if err == redis.Nil || err != nil { // redis.Nil->Key does not exists
		log.Info("Cache miss for key %s", key)
		events, err := p.Provider.Events(lat, lon, rng, sorting)

		if err != nil {
			return nil, err
		}

		p.updateCache(events, key)

		return events, nil

	} else {
		log.Info("Cache hit for key %s", key)
		var events []DojoEvent
		json.Unmarshal([]byte(cache), &events)
		return events, nil
	}

}
func (p CachedEventsProvider) updateCache(events []DojoEvent, key string) {
	eventsJson, err := json.Marshal(events)
	if err == nil {

		res := p.redis.Set(key, string(eventsJson), p.cacheDuration)
		if res.Err() != nil {
			log.Warn("Unable to Set cache for %d events with key %s error: %s\n", len(events), key, res.Err().Error())
		}else{
			log.Info("Cache update for key %s", key)
		}


	} else {
		panic(err)
		fmt.Printf("Unable to marshal %d events\n", len(events))
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