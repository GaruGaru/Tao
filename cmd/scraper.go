package cmd

import (
	"github.com/GaruGaru/Tao/providers"
	"github.com/GaruGaru/Tao/scraper"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var EventbriteToken string
var RunDelay int

func init() {

	scraperCmd.Flags().StringVarP(&EventbriteToken, "eventbrite_token", "e", "", "Eventbrite api token")

	viper.BindPFlag("eventbrite_token", scraperCmd.Flags().Lookup("eventbrite_token"))

	scraperCmd.Flags().IntVarP(&RunDelay, "scraper_delay", "d", 3600, "Scraper run delay in seconds")

	viper.BindPFlag("scraper_delay", scraperCmd.Flags().Lookup("scraper_delay"))

	rootCmd.AddCommand(scraperCmd)
}

func newEventsProvider() providers.EventProvider {
	availableProviders := make([]providers.EventProvider, 0)

	availableProviders = append(availableProviders, providers.NewZenPlatformProvider())

	if viper.GetString("eventbrite_token") != "" {
		availableProviders = append(availableProviders, providers.EventBriteProvider{ApiKey: viper.GetString("eventbrite_token")})
	}

	return providers.ProvidersManager{Providers: availableProviders}
}

func handleOsSignals(scraper *scraper.DojoScraper) {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		sig := <-c
		log.Infof("Got signal %s, terminating gracefully", sig.String())
		scraper.Lock.Release()
		log.Info("Tao scraper terminated successfully")
		os.Exit(1)
	}()
}

var scraperCmd = &cobra.Command{
	Use:   "scraper",
	Short: "Start the scraper",
	Long:  `Start the events scraper with customizable storage`,
	Run: func(cmd *cobra.Command, args []string) {

		dojoScraper := scraper.DojoScraper{
			Scraper: scraper.DefaultEventScraper{Provider: newEventsProvider()},
			Storage: GetScraperStorage(),
			Lock:    GetScraperLock(),
			Delayer: GetScraperDelayer(),
			Statter: GetStatter(),
		}

		cleaner := GetStorageCleaner()

		handleOsSignals(&dojoScraper)

		delay := uint64(viper.GetInt("scraper_delay"))

		log.WithFields(log.Fields{
			"storage": viper.GetString("storage"),
			"delay":   delay,
		}).Info("Tao scraper service started")

		gocron.Every(delay).Seconds().Do(func() {
			log.Info("Scraping started")
			err := dojoScraper.Run()
			if err != nil {
				log.Errorf("Error %s", err.Error())
			}
			log.Info("Done scraping")
		})

		gocron.Every(3).Hours().Do(func() {
			log.Info("Storage cleaner started")
			res, err := cleaner.Cleanup()
			if err != nil {
				log.Errorf("Error %s", err.Error())
			} else {
				log.Infof("Cleaner removed %d expired events", res.Removed)
			}
			log.Info("Storage cleaner scraping")
		})

		gocron.RunAll()
		<-gocron.Start()
	},
}
