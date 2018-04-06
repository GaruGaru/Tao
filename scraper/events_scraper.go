package scraper

import (
	"github.com/GaruGaru/Tao/providers"
	"github.com/go-redis/redis"
	"fmt"
	"encoding/json"
)

type EventsScraper interface {
	Scrape() ([]providers.DojoEvent, error)
	Save(events []providers.DojoEvent) error
}

type RedisEventScraper struct {
	Redis    redis.Client
	Provider providers.EventProvider
	GeoKey   string
}

func (s RedisEventScraper) Scrape() ([]providers.DojoEvent, error) {
	return s.Provider.Events(0.0, 0.0, 0, "distance")
}

func (s RedisEventScraper) Save(events []providers.DojoEvent) error {
	for _, event := range events {
		key := KeyFromEvent(event)

		jsonEvent, err := eventToJson(event)

		if err != nil {
			return err
		}

		addResult := s.Redis.Set(key, jsonEvent, 0)

		if addResult.Err() != nil {
			return addResult.Err()
		}

		geoResult := s.Redis.GeoAdd(s.GeoKey, &redis.GeoLocation{Longitude: event.Location.Longitude, Latitude: event.Location.Latitude, Name: key})

		if geoResult.Err() != nil {
			return geoResult.Err()
		}

	}

	return nil
}

func KeyFromEvent(event providers.DojoEvent) string {
	return fmt.Sprintf("%s:%d:%s", event.Title, event.StartTime, event.TicketUrl)
}

func eventToJson(events providers.DojoEvent) (string, error) {
	eventsJson, err := json.Marshal(events)
	if err != nil {
		return "", err
	}
	return string(eventsJson), nil
}
