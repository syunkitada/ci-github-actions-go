package ctl

import (
	"bytes"
	"testing"
)

func TestStart(t *testing.T) {
	cmd := RootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"start"})
	cmd.Execute()
}
