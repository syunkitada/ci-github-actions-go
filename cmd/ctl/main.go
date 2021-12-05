package main

import (
	"github.com/syunkitada/ci-github-actions-go/pkg/ctl"
	"github.com/syunkitada/ci-github-actions-go/pkg/lib/logger"
)

func main() {
	if err := ctl.RootCmd.Execute(); err != nil {
		logger.StdoutFatal(err)
	}
}
