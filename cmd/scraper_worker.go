package cmd

import (
	"github.com/spf13/cobra"
	worker "github.com/contribsys/faktory_worker_go"
	"time"
	"github.com/sirupsen/logrus"
)

func init() {

	rootCmd.AddCommand(scraperWorkerCmd)
}

var scraperWorkerCmd = &cobra.Command{
	Use:   "scraper-worker",
	Short: "Start the scraper worker",
	Long:  `Start the events scraper worker with customizable storage`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr := worker.NewManager() // FAKTORY_PROVIDER=tcp://faktory.example.com:12345

		mgr.Register("probe", probe)
		mgr.Concurrency = 10
		mgr.Queues = []string{"scraping", "default"}

		mgr.On(worker.Shutdown, func() {
			// TODO Dispose here resources
		})

		mgr.Run()
	},
}

func probe(ctx worker.Context, args ...interface{}) error {
	time.Sleep(1 * time.Second)
	logrus.Info("Async task probe with %d args", len(args))
	return nil
}
