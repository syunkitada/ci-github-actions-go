package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("DEBUG stop")
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
