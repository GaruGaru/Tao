package cmd

import (
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/GaruGaru/Tao/scraper"
	"github.com/GaruGaru/Tao/providers"
	log "github.com/sirupsen/logrus"
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

	return providers.ProvidersManager{Providers: availableProviders,}
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
			Statter: GetStatter(),
		}

		delayer := GetScraperDelayer()

		delay := uint64(viper.GetInt("scraper_delay"))

		log.WithFields(log.Fields{
			"storage": viper.GetString("storage"),
			"delay":   delay,
		}).Info("Tao scraper service started")

		gocron.Every(delay).Seconds().Do(func() {
			canRun, err := delayer.CanRun()

			if err != nil {
				log.Error(err.Error())
			} else {
				if canRun {
					log.Info("Scraping started")
					err := dojoScraper.Run()

					if err != nil {
						log.Errorf("Error %s", err.Error())
					}
					delayer.Refresh()
					log.Info("Done scraping")
				}else{
					log.Infof("Can't run scraper, tried too early")
				}
			}

		})

		gocron.RunAll()
		<-gocron.Start()

	},
}
