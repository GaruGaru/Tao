package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/GaruGaru/Tao/providers"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type EventsStorage interface {
	Store(events []providers.DojoEvent) error
}

type InMemoryEventsStorage struct {
	Storage map[string]string
}

func NewInMemoryEventsStorage() InMemoryEventsStorage {
	return InMemoryEventsStorage{
		Storage: make(map[string]string),
	}
}

func (m InMemoryEventsStorage) Store(events []providers.DojoEvent) error {
	for _, e := range events {
		jsonEvent, err := eventToJson(e)
		if err != nil {
			return err
		}
		m.Storage[keyFromEvent(e)] = jsonEvent
	}
	return nil
}

type FileSystemEventsStorage struct {
	StoreFile string
}

func (m FileSystemEventsStorage) Store(events []providers.DojoEvent) error {
	content, err := json.Marshal(events)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(m.StoreFile, content, 0644)
	return err
}

type RedisEventsStorage struct {
	Redis  redis.Client
	GeoKey string
}

func (p RedisEventsStorage) Store(events []providers.DojoEvent) error {
	for _, event := range events {

		key := keyFromEvent(event)

		jsonEvent, err := eventToJson(event)

		if err != nil {
			return err
		}

		eventStart := time.Unix(event.StartTime, 0)

		expiration := eventStart.Sub(time.Now())

		if expiration > 0 {
			addResult := p.Redis.Set(key, jsonEvent, expiration)

			if addResult.Err() != nil {
				return addResult.Err()
			}

			geoResult := p.Redis.GeoAdd(p.GeoKey, &redis.GeoLocation{Longitude: event.Location.Longitude, Latitude: event.Location.Latitude, Name: key})

			if geoResult.Err() != nil {
				return geoResult.Err()
			}
		}

	}

	saveResult := p.Redis.Save()
	if saveResult.Err() != nil {
		logrus.Warn("Redis auto save is running, skipped Save command")
	}

	return nil
}

func keyFromEvent(event providers.DojoEvent) string {
	return fmt.Sprintf("%s:%d:%s", event.Title, event.StartTime, event.TicketUrl)
}

func eventToJson(events providers.DojoEvent) (string, error) {
	eventsJson, err := json.Marshal(events)
	if err != nil {
		return "", err
	}
	return string(eventsJson), nil
}
