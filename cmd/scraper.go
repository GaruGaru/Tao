package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/go-redis/redis"
	"github.com/GaruGaru/Tao/scraper"
	"github.com/GaruGaru/Tao/providers"
)

var EventbriteToken string

func init() {

	scraperCmd.Flags().StringVarP(&EventbriteToken, "eventbrite_token", "e", "", "Eventbrite api token")

	viper.BindPFlag("eventbrite_token", scraperCmd.Flags().Lookup("eventbrite_token"))

	rootCmd.AddCommand(scraperCmd)
}

var scraperCmd = &cobra.Command{
	Use:   "scraper",
	Short: "Start the scraper",
	Long:  `Start the events scraper with customizable storage`,
	Run: func(cmd *cobra.Command, args []string) {

		redisClient := redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis_host"),
			DB:       0,
		})

		eventsProvider := providers.ProvidersManager{
			Providers: []providers.EventProvider{providers.EventBrite{ApiKey: viper.GetString("eventbrite_token")}},
		}

		dojoScraper := scraper.DojoScraper{
			Scraper: scraper.DefaultEventScraper{Provider: eventsProvider},
			Storage: scraper.RedisEventsStorage{Redis: *redisClient, GeoKey: "locations"},
			Lock:    scraper.FileSystemLock{LockFile: "/tmp/tao.lock"},
		}

		err := dojoScraper.Run()

		if err != nil {
			panic(err)
		}

	},
}
