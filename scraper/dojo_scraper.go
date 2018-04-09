package scraper

import (
	"fmt"
)

type DojoScraper struct {
	Storage EventsStorage
	Scraper EventsScraper
	Lock    DojoScraperLock
}

func (d DojoScraper) Run() error {

	if d.Lock.Obtain() == nil {
		defer d.Lock.Release()

		fmt.Println("Starting scraper")
		events, err := d.Scraper.Scrape()
		fmt.Println("Done scraper")
		if err != nil {
			return err
		}

		err = d.Storage.Store(events)

		return err

	} else {
		fmt.Printf("Scraper already running")
		return nil
	}

}
