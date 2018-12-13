package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(printargsCmd)
}

var printargsCmd = &cobra.Command{
	Use: "printargs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("args: %s\n", strings.Join(args, " "))
	},
}
