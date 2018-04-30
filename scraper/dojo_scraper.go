package scraper

import (
	"github.com/cactus/go-statsd-client/statsd"
	"time"
	log "github.com/sirupsen/logrus"
	"fmt"
)

type DojoScraper struct {
	Storage EventsStorage
	Scraper EventsScraper
	Lock    DojoScraperLock
	Delayer ScraperDelayer
	Statter statsd.Statter
}

func (d DojoScraper) Run() error {

	canRun, err := d.Delayer.CanRun()

	if err != nil {
		return fmt.Errorf("delayer error: %s", err.Error())
	}

	if !canRun {
		log.Info("Run delayed, tried too early")
		return nil
	}

	err = d.Lock.Obtain()

	if err == nil {

		defer d.Delayer.Refresh()
		defer d.Lock.Release()

		log.Info("Starting scraper")
		d.Statter.Inc("scraper.run", 1, 1)

		startScrape := time.Now()

		events, err := d.Scraper.Scrape()

		log.Info("Done scraper")
		if err != nil {
			d.Statter.Inc("scraper.scraping.error", 1, 1)
			return err
		}else{
			d.Statter.TimingDuration("scraper.scraping.latency", time.Now().Sub(startScrape), 1.0)
		}

		startStore := time.Now()

		err = d.Storage.Store(events)

		if err != nil {
			d.Statter.Inc("scraper.storage.error", 1, 1)
		}else{
			d.Statter.TimingDuration("scraper.storage.latency", time.Now().Sub(startStore), 1.0)
		}

		return err

	} else {
		log.Infof("Scraper already running %s", err.Error())
		d.Statter.Inc("scraper.already_running", 1, 1)
		return nil
	}

}
