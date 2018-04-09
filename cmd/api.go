package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/go-redis/redis"
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

		statsdClient, err := statsd.NewClient(viper.GetString("statsd_host"), "tao")

		if err != nil {
			panic(err)
		}

		redisClient := redis.NewClient(&redis.Options{
			Addr: viper.GetString("redis_host"),
			Password: "",
			DB:       0,
		})

		localCacheExpiration := 15 * time.Minute
		remoteCacheExpiration := 30 * time.Minute

		caches := []providers.EventsCache{
			providers.NewLocalCache(localCacheExpiration),
			providers.NewRedisEventsCache(*redisClient, remoteCacheExpiration),
		}

		eventsProvider := providers.ProvidersManager{
			Providers: []providers.EventProvider{providers.NewEventBriteProvider()},
		}

		cachedProvider := providers.NewCachedEventsProvider(eventsProvider, caches, statsdClient)

		taoApi := api.EventsApi{Provider: cachedProvider, Statsd: statsdClient}

		taoApi.Run(viper.GetInt("port"))
	},
}
