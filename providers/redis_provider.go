package providers

import (
	"github.com/go-redis/redis"
	"encoding/json"
	"fmt"
	"sync"
)

type RedisEventsProvider struct {
	Redis        redis.Client
	LocationsKey string
}

func (r RedisEventsProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {

	locations := r.Redis.GeoRadius(r.LocationsKey, lon, lat, &redis.GeoRadiusQuery{
		Radius: float64(rng),
	})

	if locations.Err() != nil && locations.Err() != redis.Nil {
		return nil, locations.Err()
	}

	results := locations.Val()

	var wg sync.WaitGroup

	eventsChannel := make(chan DojoEvent, len(results))

	wg.Add(len(results))

	for _, l := range results {
		go fetchEventInfo(l, r, eventsChannel, &wg)
	}

	wg.Wait()
	close(eventsChannel)

	var dojoEvents []DojoEvent

	for e := range eventsChannel {
		dojoEvents = append(dojoEvents, e)
	}

	return dojoEvents, nil
}
func fetchEventInfo(l redis.GeoLocation, r RedisEventsProvider, eventsChannel chan DojoEvent, wg *sync.WaitGroup) {
	defer wg.Done()
	key := l.Name
	get := r.Redis.Get(key)
	event, err := eventFromJson(get.Val())
	if get.Err() != nil || err != nil {
		fmt.Printf("Unable to get event from json for key %s\n", key)
	} else {
		eventsChannel <- event
	}
}

func eventFromJson(data string) (DojoEvent, error) {
	var event DojoEvent
	err := json.Unmarshal([]byte(data), &event)
	if err != nil {
		return DojoEvent{}, err
	} else {
		return event, nil
	}
}
