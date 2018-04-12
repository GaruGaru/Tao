package cmd

import (
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var version = "1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Tao",
	Long:  `Print the verbose version number of Tao`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Tao version %s", version)
	},
}
