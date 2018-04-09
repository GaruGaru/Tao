package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/GaruGaru/Tao/scraper"
	"github.com/GaruGaru/Tao/providers"
	"fmt"
	"time"
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

	if viper.GetString("eventbrite_token") != "" {
		availableProviders = append(availableProviders, providers.EventBrite{ApiKey: viper.GetString("eventbrite_token")})
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
		}

		delayer := GetScraperDelayer()

		for ; ; {
			canRun, err := delayer.CanRun()

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			if canRun {
				err := dojoScraper.Run()

				if err != nil {
					fmt.Printf("Scraping failed: %s\n", err.Error())
				}

				time.Sleep(time.Duration(viper.GetInt("scraper_delay")) * time.Second)
				delayer.Refresh()
			}

		}

	},
}
