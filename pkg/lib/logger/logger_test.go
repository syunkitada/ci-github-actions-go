package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	a := assert.New(t)

	{
		Init(&Config{EnableTest: true, LogDir: "/tmp", Host: "host", App: "app"})
		tctx := NewTraceContext(&NewTraceContextInput{Host: "hoge", App: "app"})
		Info(tctx, "test")

		Init(&Config{EnableTest: true, LogDir: "/notfound"})
	}

	Init(&Config{EnableTest: true, LogDir: "stdout"})
	tctx := NewTraceContext(&NewTraceContextInput{})

	StdoutInfo("test")
	StdoutInfof("test %s", "test")

	StdoutFatal("test")
	StdoutFatalf("test %s", "test")

	Fatal(tctx, "test")
	Fatalf(tctx, "test %s", "test")

	Info(tctx, "test")
	Infof(tctx, "test %s", "test")

	Warn(tctx, "test")
	Warnf(tctx, "test %s", "test")

	Error(tctx, "test")
	Errorf(tctx, "test %s", "test")

	startTime := StartTrace(tctx)
	EndTrace(tctx, startTime, nil, 0)
	EndTrace(tctx, startTime, fmt.Errorf("error"), 0)
	EndTrace(tctx, startTime, nil, 10)

	tid := tctx.GetTraceId()
	fmt.Println(tid)

	tctx.SetMetadata(map[string]string{
		"hoge": "hoge",
	})
	Info(tctx, "test")
	tctx.ResetMetadata()

	a.Equal(123, 123, "they should be equal")
}
