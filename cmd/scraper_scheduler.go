package cmd

import (
	"github.com/spf13/cobra"
	faktory "github.com/contribsys/faktory/client"
	"fmt"
)

func init() {
	rootCmd.AddCommand(scraperSchedulerCmd)
}

var scraperSchedulerCmd = &cobra.Command{
	Use:   "scraper-scheduler",
	Short: "Start the scraper scheduler",
	Long:  `Start the events scraper scheduler with customizable storage`,
	Run: func(cmd *cobra.Command, args []string) {

		cl, err := faktory.Open() // FAKTORY_PROVIDER=tcp://faktory.example.com:12345
		if err != nil {
			panic(err)
		}

		job := faktory.NewJob("SomeJob", 1, 2, "3")

		err = cl.Push(job)

		if err != nil {
			panic(err)
		}

		fmt.Println(cl.Info())

	},
}

