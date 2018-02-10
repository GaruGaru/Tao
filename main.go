package main

import (
	"github.com/GaruGaru/Tao/providers"
	"github.com/GaruGaru/Tao/api"
	"github.com/smira/go-statsd"
	"os"
)

func main() {

	statsdClient := statsd.NewClient(getenv("STATSD_HOST", "localhost:8125"))

	provider := providers.NewCachedEventsProvider(providers.ProvidersManager{
		Providers: []providers.EventProvider{
			providers.NewEventBriteProvider(),
		},
	}, *statsdClient)

	taoApi := api.EventsApi{Provider: provider}

	taoApi.Run()

}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
