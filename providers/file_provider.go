package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func readFile(file string) ([]DojoEvent, error) {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %s", file, err.Error())
	}

	var events []DojoEvent

	err = json.Unmarshal(content, &events)

	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal json: %s  %s", content, err.Error())
	}

	return events, nil
}

func NewFileEventProvider(file string) FileEventProvider {
	events, err := readFile(file)
	if err != nil {
		panic(err)
	}
	return FileEventProvider{events: events}
}

type FileEventProvider struct {
	events []DojoEvent
}

func (e FileEventProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {
	return e.events, nil
}
