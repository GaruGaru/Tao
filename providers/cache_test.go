package providers

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"time"
)

func loadTestData(t *testing.T) []DojoEvent {
	testJson, err := ioutil.ReadFile("testdata/provider_response.json")

	if err != nil {
		t.Log("Unable to read test json file")
		t.FailNow()
	}

	var events []DojoEvent

	err = json.Unmarshal(testJson, &events)

	if err != nil {
		t.Log("Unable to unmarshal test json file")
		t.FailNow()
	}

	return events
}

func TestLocalCache(t *testing.T) {

	testEvents := loadTestData(t)

	if len(testEvents) == 0 {
		t.FailNow()
	}

	cache := NewLocalCache(1 * time.Minute)

	key := RequestKey(1.0, 1.0, 1, "distance")
	cache.Put(key, testEvents)

	cachedEvents, present, _ := cache.Get(1.0, 1.0, 1, "distance")

	if !present {
		t.Log("Cache miss for " + key)
		t.FailNow()
	}

	if len(cachedEvents) == 0 {
		t.Log("Event cache hit but no events returned")
		t.FailNow()
	}

	if len(cachedEvents) != len(testEvents) {
		t.Log("Event cache hit but different event count")
		t.FailNow()
	}

}
