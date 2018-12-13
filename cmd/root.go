package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cobra_example",
	Short: "cobra_example is am simple example for the usage of cobra library.",
	Long: `An example for the usage of cobra library with the poor translate of README.md of
	            github.com/spf13/cobra
                sourcecode is available at http://github.com/xiaotuanyu120/cobra_example`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("run cobra_example success!\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
