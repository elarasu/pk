package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version info of pf",
	Long:  `All software has versions. This is pf's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pf v0.1 -- HEAD")
	},
}
