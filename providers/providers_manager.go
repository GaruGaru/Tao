package providers

import (
	"sort"
	"sync"
	"github.com/sirupsen/logrus"
)

type ProvidersManager struct {
	Providers []EventProvider
}

func (m ProvidersManager) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {

	var wg sync.WaitGroup

	providersCount := len(m.Providers)

	eventsChannel := make(chan []DojoEvent, providersCount)
	wg.Add(providersCount)

	for _, provider := range m.Providers {
		go fetchEvents(provider, lat, lon, rng, sorting, eventsChannel, &wg)
	}

	wg.Wait()

	close(eventsChannel)

	var dojoEvents []DojoEvent

	for providerEvents := range eventsChannel {
		for _, event := range providerEvents {
			dojoEvents = append(dojoEvents, event)
		}
	}

	if sorting == "distance" {
		sort.Slice(dojoEvents, func(i, j int) bool { return dojoEvents[i].Location.Distance < dojoEvents[j].Location.Distance })
	} else if sorting == "date" {
		sort.Slice(dojoEvents, func(i, j int) bool { return dojoEvents[i].StartTime < dojoEvents[j].StartTime })
	}

	if dojoEvents == nil {
		dojoEvents = []DojoEvent{}
	}

	return dojoEvents, nil

}
func fetchEvents(provider EventProvider, lat float64, lon float64, rng int, sorting string, eventsChannel chan []DojoEvent, wg *sync.WaitGroup) {
	defer wg.Done()
	events, err := provider.Events(lat, lon, rng, sorting)
	if err == nil {
		eventsChannel <- events
	}else{
		logrus.Error(err.Error())
	}
}
