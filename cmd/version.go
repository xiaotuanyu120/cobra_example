package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var Version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version of cobra_example",
	Long:  "All software has version, here is cobra_example's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cobra_example is on version %s\n", Version)
	},
}
