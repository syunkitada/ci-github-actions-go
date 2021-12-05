package ctl

import (
	"github.com/spf13/cobra"

	"github.com/syunkitada/ci-github-actions-go/pkg/lib/logger"
)

var RootCmd = &cobra.Command{}

func init() {
	logger.Init(&logger.Config{LogDir: "stdout"})
}
