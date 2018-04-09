package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
	"github.com/GaruGaru/Tao/providers"
	"github.com/GaruGaru/Tao/api"
)

var Port int

func init() {

	apiCmd.Flags().IntVarP(&Port, "port", "p", 8080, "Api server port")
	viper.BindPFlag("port", apiCmd.Flags().Lookup("port"))
	viper.AutomaticEnv()
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "serve-api",
	Short: "Start the api server",
	Long:  `Start the events api server`,
	Run: func(cmd *cobra.Command, args []string) {

		statter := GetStatsdClient()

		localCacheExpiration := 15 * time.Minute
		remoteCacheExpiration := 30 * time.Minute

		caches := []providers.EventsCache{
			providers.NewLocalCache(localCacheExpiration),
			providers.NewRedisEventsCache(*GetRedisClient(), remoteCacheExpiration),
		}

		eventsProvider := providers.ProvidersManager{
			Providers: []providers.EventProvider{providers.NewEventBriteProvider()},
		}

		cachedProvider := providers.NewCachedEventsProvider(eventsProvider, caches, statter)

		taoApi := api.EventsApi{Provider: cachedProvider, Statsd: statter}

		taoApi.Run(viper.GetInt("port"))
	},
}
