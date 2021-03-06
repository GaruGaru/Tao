package providers

import (
	"encoding/json"
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/go-redis/redis"
	"sort"
	"sync"
)

type RedisEventsProvider struct {
	Redis        redis.Client
	LocationsKey string
	Statter      statsd.Statter
}

func (r RedisEventsProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {

	fmt.Printf("Geo query with lat: %f, lon: %f, range: %f", lat, lon, float64(rng))

	locations := r.Redis.GeoRadius(r.LocationsKey, lon, lat, &redis.GeoRadiusQuery{
		Radius: float64(rng),
		Unit:   "km",
	})

	if locations.Err() != nil && locations.Err() != redis.Nil {
		return nil, locations.Err()
	}

	results := locations.Val()

	var wg sync.WaitGroup

	eventsChannel := make(chan DojoEvent, len(results))

	wg.Add(len(results))

	for _, geoLocation := range results {
		go fetchEventInfo(lat, lon, geoLocation, r, eventsChannel, &wg)
	}

	wg.Wait()
	close(eventsChannel)

	var dojoEvents = make([]DojoEvent, 0)

	for e := range eventsChannel {
		dojoEvents = append(dojoEvents, e)
	}

	if sorting == "distance" {
		sort.Slice(dojoEvents, func(i, j int) bool { return dojoEvents[i].Location.Distance < dojoEvents[j].Location.Distance })
	} else if sorting == "date" {
		sort.Slice(dojoEvents, func(i, j int) bool { return dojoEvents[i].StartTime < dojoEvents[j].StartTime })
	}

	return dojoEvents, nil
}

func fetchEventInfo(hLat float64, hLon float64, geoLocation redis.GeoLocation, r RedisEventsProvider, eventsChannel chan DojoEvent, wg *sync.WaitGroup) {
	defer wg.Done()
	key := geoLocation.Name
	get := r.Redis.Get(key)
	event, err := eventFromJson(get.Val())
	if get.Err() != nil || err != nil {
		fmt.Printf("Unable to get event from json for key %s\n", key)
	} else {
		event.Location.Distance = Distance(hLat, hLon, event.Location.Latitude, event.Location.Longitude)
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
