package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("DEBUG start")
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
