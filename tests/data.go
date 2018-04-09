package tests

import (
	"testing"
	"github.com/GaruGaru/Tao/providers"
	"io/ioutil"
	"encoding/json"
)

func LoadTestEvents(t *testing.T, path string) []providers.DojoEvent {
	testJson, err := ioutil.ReadFile(path)

	if err != nil {
		t.Log("Unable to read test json file")
		t.FailNow()
	}

	var events []providers.DojoEvent

	err = json.Unmarshal(testJson, &events)

	if err != nil {
		t.Log("Unable to unmarshal test json file")
		t.FailNow()
	}

	return events
}
