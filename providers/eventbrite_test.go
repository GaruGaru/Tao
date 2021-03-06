package providers

import (
	"fmt"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"path"
	"testing"
)

func readTestDataRaw(t *testing.T, file string) string {
	rawContent, err := ioutil.ReadFile(path.Join("testdata", file))
	if err != nil {
		t.Log("Unable to read test json file")
		t.FailNow()
	}
	return string(rawContent)
}

func TestFetchArticlesWithoutPagination(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://www.eventbriteapi.com").
		Get("/v3/events/search/").
		Reply(200).
		BodyString(readTestDataRaw(t, "eventbrite_response_no_pagination.json"))

	gock.New("https://www.eventbriteapi.com").
		Get("/v3/venues/(.*)").
		Persist().
		Reply(200).
		BodyString(readTestDataRaw(t, "eventbrite_venue_response.json"))

	eventbrite := EventBriteProvider{ApiKey: "test"}

	events, err := eventbrite.Events(10, 10, 1000, "distance")

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(events) != 50 {
		t.Log("Fetched events count does not match != 50")
		t.FailNow()
	}

}

func TestFetchArticlesWithPagination(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://www.eventbriteapi.com").
		Get("/v3/events/search/").
		Persist().
		Reply(200).
		BodyString(readTestDataRaw(t, "eventbrite_response_p1.json"))

	gock.New("https://www.eventbriteapi.com").
		Get("/v3/events/search/").
		MatchParam("page", "2").
		Reply(200).
		BodyString(readTestDataRaw(t, "eventbrite_response_p2.json"))

	gock.New("https://www.eventbriteapi.com").
		MatchParam("page", "3").
		Reply(200).
		BodyString(readTestDataRaw(t, "eventbrite_response_p3.json"))

	gock.New("https://www.eventbriteapi.com").
		Get("/v3/venues/(.*)").
		Persist().
		Reply(200).
		BodyString(readTestDataRaw(t, "eventbrite_venue_response.json"))

	eventbrite := EventBriteProvider{ApiKey: "test"}

	events, err := eventbrite.Events(10, 10, 1000, "distance")

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(events) != 150 {
		t.Log(fmt.Sprintf("Fetched events count does not match  %d != 150", len(events)))
		t.FailNow()
	}

}

func TestFetchArticlesWithPaginationWithRateLimit(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.Off()
	gock.New("https://www.eventbriteapi.com").
		Get("/v3/events/search/").
		MatchParam("page", "1").
		Persist().
		Reply(200).
		BodyString(readTestDataRaw(t, "eventbrite_response_p1.json"))

	gock.New("https://www.eventbriteapi.com").
		Get("/v3/events/search/").
		MatchParam("page", "2").
		Reply(403).
		BodyString("Rate-limited")

	gock.New("https://www.eventbriteapi.com").
		MatchParam("page", "3").
		Reply(403).
		BodyString("Rate-limited")

	gock.New("https://www.eventbriteapi.com").
		Get("/v3/venues/(.*)").
		Persist().
		Reply(200).
		BodyString(readTestDataRaw(t, "eventbrite_venue_response.json"))

	eventbrite := EventBriteProvider{ApiKey: "test"}

	events, err := eventbrite.Events(10, 10, 1000, "distance")

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(events) != 50 {
		t.Log(fmt.Sprintf("Fetched events count does not match  %d != 50", len(events)))
		t.FailNow()
	}

}
