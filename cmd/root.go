package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/spf13/viper"
)

var Storage string

var RedisHost string

var StatsdHost string

func init() {

	rootCmd.PersistentFlags().StringVarP(&Storage, "storage", "s", "memory", "Values storage")
	viper.BindPFlag("storage", rootCmd.PersistentFlags().Lookup("storage"))

	rootCmd.PersistentFlags().StringVarP(&RedisHost, "redis_host", "r", "localhost:6379", "Redis storage host")
	viper.BindPFlag("redis_host", rootCmd.PersistentFlags().Lookup("redis_host"))

	rootCmd.PersistentFlags().StringVarP(&StatsdHost, "statsd_host", "", "localhost:8125", "Statsd metrics host")
	viper.BindPFlag("statsd_host", rootCmd.PersistentFlags().Lookup("statsd_host"))

	viper.AutomaticEnv()

}



var rootCmd = &cobra.Command{
	Use:   "tao",
	Short: "Tao is an aggregator for coderdojo events",
	Long: `Tao is a fast, scalable, container ready aggregator for coderdojo events`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Tao is a fast, scalable, container ready aggregator for coderdojo events")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

