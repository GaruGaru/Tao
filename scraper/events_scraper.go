package scraper

import (
	"github.com/GaruGaru/Tao/providers"
)

type EventsScraper interface {
	Scrape() ([]providers.DojoEvent, error)
}

type DefaultEventScraper struct {
	Provider providers.EventProvider
}

func (s DefaultEventScraper) Scrape() ([]providers.DojoEvent, error) {
	return s.Provider.Events(0.0, 0.0, 1<<63 - 1, "distance")
}
