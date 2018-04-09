package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
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
		fmt.Printf("Tao • coderdojo events aggregator %s", version)
	},
}
