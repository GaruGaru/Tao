package scraper

import (
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"time"
	log "github.com/sirupsen/logrus"
)

type DojoScraper struct {
	Storage EventsStorage
	Scraper EventsScraper
	Lock    DojoScraperLock
	Statter statsd.Statter
}

func (d DojoScraper) Run() error {

	err := d.Lock.Obtain()
	if err == nil {
		fmt.Println("Starting scraper")
		d.Statter.Inc("scraper.run", 1, 1)

		startScrape := time.Now()
		events, err := d.Scraper.Scrape()
		d.Statter.TimingDuration("scraper.scraping.latency", time.Now().Sub(startScrape), 1)

		fmt.Println("Done scraper")
		if err != nil {
			d.Statter.Inc("scraper.scraping.error", 1, 1)
			return err
		}

		startStore := time.Now()
		err = d.Storage.Store(events)

		if err != nil{
			d.Statter.Inc("scraper.storage.error", 1, 1)
		}
		d.Statter.TimingDuration("scraper.storage.latency", time.Now().Sub(startStore), 1)

		releaseErr := d.Lock.Release()
		if releaseErr != nil {
			return releaseErr
		}

		return err

	} else {
		fmt.Printf("Scraper already running: %s\n", err)
		d.Statter.Inc("scraper.already_running", 1, 1)
		return nil
	}

}
