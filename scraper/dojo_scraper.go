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


		fmt.Println("Starting scraper")
		events, err := d.Scraper.Scrape()
		fmt.Println("Done scraper")
		if err != nil {
			return err
		}

		err = d.Storage.Store(events)
		releaseErr := d.Lock.Release()
		if releaseErr != nil {
			return releaseErr
		}
		return err

	} else {
		fmt.Printf("Scraper already running")
		return nil
	}

}
