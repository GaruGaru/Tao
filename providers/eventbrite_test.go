package providers

import (
	"testing"
	"io/ioutil"
	"path"
	"gopkg.in/h2non/gock.v1"
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

	eventbrite := EventBrite{ApiKey: "test"}

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

	eventbrite := EventBrite{ApiKey: "test"}

	events, err := eventbrite.Events(10, 10, 1000, "distance")

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(events) != 150 {
		t.Log("Fetched events count does not match != 150")
		t.FailNow()
	}

}
