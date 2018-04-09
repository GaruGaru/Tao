package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/GaruGaru/Tao/scraper"
	"github.com/GaruGaru/Tao/providers"
)

var EventbriteToken string

func init() {

	scraperCmd.Flags().StringVarP(&EventbriteToken, "eventbrite_token", "e", "", "Eventbrite api token")

	viper.BindPFlag("eventbrite_token", scraperCmd.Flags().Lookup("eventbrite_token"))

	rootCmd.AddCommand(scraperCmd)
}

func newEventsProvider() providers.EventProvider {
	availableProviders := make([]providers.EventProvider, 1)

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

		err := dojoScraper.Run()

		if err != nil {
			panic(err)
		}

	},
}
