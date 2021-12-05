package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/syunkitada/ci-github-actions-go/pkg_infra/lib/os_utils"
)

var (
	name          string
	host          string
	app           string
	Logger        *log.Logger
	stdoutLogger  *log.Logger
	LogTimeFormat string
	enableTest    bool
)

const (
	infoLog    = "INFO"
	warningLog = "WARNING"
	errorLog   = "ERROR"
	fatalLog   = "FATAL"
)

type TraceContext struct {
	mtx      *sync.Mutex
	code     uint8
	err      error
	host     string
	app      string
	function string
	traceId  string
	metadata map[string]string
}

func (tctx *TraceContext) GetTraceId() string {
	return tctx.traceId
}

func (tctx *TraceContext) SetMetadata(data map[string]string) {
	tctx.mtx.Lock()
	for k, v := range data {
		tctx.metadata[k] = v
	}
	tctx.mtx.Unlock()
}

func (tctx *TraceContext) ResetMetadata() {
	tctx.mtx.Lock()
	tctx.metadata = map[string]string{}
	tctx.mtx.Unlock()
}

type NewTraceContextInput struct {
	Host string
	App  string
}

func NewTraceContext(input *NewTraceContextInput) (tctx *TraceContext) {
	tctx = &TraceContext{
		mtx:      new(sync.Mutex),
		traceId:  xid.New().String(),
		metadata: map[string]string{},
		code:     0,
		err:      nil,
	}
	if input.Host == "" {
		tctx.host = host
	} else {
		tctx.host = input.Host
	}
	if input.App == "" {
		tctx.app = app
	} else {
		tctx.app = input.App
	}
	return
}

type Config struct {
	LogTimeFormat string
	LogDir        string
	Host          string
	App           string
	EnableTest    bool
}

func Init(conf *Config) {
	stdoutLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	LogTimeFormat = conf.LogTimeFormat
	enableTest = conf.EnableTest
	host = conf.Host
	app = conf.App

	name = os.Getenv("LOG_FILE")
	if name == "" {
		for _, arg := range os.Args {
			option := strings.Index(arg, "-")
			if option == 0 {
				continue
			}
			slash := strings.LastIndex(arg, "/")
			if slash > 0 {
				arg = arg[slash+1:]
			}
			name += "-" + arg
		}
		name = name[1:]
	}

	if conf.LogDir == "stdout" {
		Logger = log.New(os.Stdout, "", 0)
	} else {
		logfilePath := filepath.Join(conf.LogDir, name+".log")
		logfile, err := os.OpenFile(logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			StdoutFatalf("Failed open logfile: %v", logfile)
		} else {
			Logger = log.New(logfile, "", 0)
		}
	}
}

func timePrefix() string {
	return "Time=\"" + time.Now().Format(LogTimeFormat) + "\""
}

func convertTags(tctx *TraceContext) string {
	tags := []string{"TraceId=\"" + tctx.traceId + "\""}
	if tctx.host != "" {
		tags = append(tags, "Host=\""+tctx.host+"\"")
	}
	if tctx.app != "" {
		tags = append(tags, "App=\""+tctx.app+"\"")
	}
	tags = append(tags, "Func=\""+tctx.function+"\"")
	for k, v := range tctx.metadata {
		tags = append(tags, k+"=\""+v+"\"")
	}
	return strings.Join(tags, " ")
}

func getFunc(depth int) string {
	var function string
	pc, file, line, ok := runtime.Caller(2 + depth)
	if !ok {
		file = ""
		line = 1
		function = ""
	} else {
		slash := strings.LastIndex(file, "/")
		if slash > 0 {
			file = file[slash+1:]
		}

		function = runtime.FuncForPC(pc).Name()
		dot := strings.LastIndex(function, ".")
		if dot > 0 {
			function = function[dot+1:]
		}
	}
	return fmt.Sprintf("%s:%d:%s", file, line, function)
}

func StdoutInfo(format string, args ...interface{}) {
	stdoutLogger.Print(infoLog + " " + fmt.Sprint(args...))
}

func StdoutInfof(format string, args ...interface{}) {
	stdoutLogger.Print(infoLog + " " + fmt.Sprintf(format, args...))
}

func StdoutFatal(args ...interface{}) {
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	l.Print(fatalLog + " " + fmt.Sprint(args...))
	os_utils.Exit(1, enableTest)
}

func StdoutFatalf(format string, args ...interface{}) {
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	l.Print(fatalLog + " " + fmt.Sprintf(format, args...))
	os_utils.Exit(1, enableTest)
}

func Fatal(tctx *TraceContext, args ...interface{}) {
	tctx.mtx.Lock()
	tctx.function = getFunc(0)
	Logger.Print(timePrefix() + " Level=\"" + fatalLog +
		"\" Msg=\"" + fmt.Sprint(args...) + "\"" + convertTags(tctx))
	tctx.mtx.Unlock()
	os_utils.Exit(1, enableTest)
}

func Fatalf(tctx *TraceContext, format string, args ...interface{}) {
	tctx.mtx.Lock()
	tctx.function = getFunc(0)
	Logger.Print(timePrefix() + " Level=\"" + fatalLog +
		"\" Msg=\"" + fmt.Sprintf(format, args...) + "\"" + convertTags(tctx))
	tctx.mtx.Unlock()
	os_utils.Exit(1, enableTest)
}

func Info(tctx *TraceContext, args ...interface{}) {
	tctx.mtx.Lock()
	tctx.function = getFunc(0)
	Logger.Print(timePrefix() + " Level=\"" + infoLog +
		"\" Msg=\"" + fmt.Sprint(args...) + "\"" + convertTags(tctx))
	tctx.mtx.Unlock()
}

func Infof(tctx *TraceContext, format string, args ...interface{}) {
	tctx.mtx.Lock()
	tctx.function = getFunc(0)
	Logger.Print(timePrefix() + " Level=\"" + infoLog +
		"\" Msg=\"" + fmt.Sprintf(format, args...) + "\"" + convertTags(tctx))
	tctx.mtx.Unlock()
}

func Warn(tctx *TraceContext, args ...interface{}) {
	tctx.mtx.Lock()
	tctx.function = getFunc(0)
	Logger.Print(timePrefix() + " Level=\"" + warningLog +
		"\" Msg=\"" + fmt.Sprint(args...) + "\"" + convertTags(tctx))
	tctx.mtx.Unlock()
}

func Warnf(tctx *TraceContext, format string, args ...interface{}) {
	tctx.mtx.Lock()
	tctx.function = getFunc(0)
	Logger.Print(timePrefix() + " Level=\"" + warningLog +
		"\" Msg=\"" + fmt.Sprintf(format, args...) + "\"" + convertTags(tctx))
	tctx.mtx.Unlock()
}

func Error(tctx *TraceContext, args ...interface{}) {
	tctx.mtx.Lock()
	tctx.function = getFunc(0)
	Logger.Print(timePrefix() + " Level=\"" + errorLog +
		"\" Msg=\"" + fmt.Sprint(args...) + "\"" + convertTags(tctx))
	tctx.mtx.Unlock()
}

func Errorf(tctx *TraceContext, format string, args ...interface{}) {
	tctx.mtx.Lock()
	tctx.function = getFunc(0)
	Logger.Print(timePrefix() + " Level=\"" + errorLog +
		"\" Msg=\"" + fmt.Sprintf(format, args...) + "\"" + convertTags(tctx))
	tctx.mtx.Unlock()
}

func StartTrace(tctx *TraceContext) time.Time {
	startTime := time.Now()
	tctx.mtx.Lock()
	tctx.function = getFunc(0)
	Logger.Print(timePrefix() + " Level=\"" + infoLog + "\" Msg=\"StartTrace\"" + convertTags(tctx))
	tctx.mtx.Unlock()
	return startTime
}

func EndTrace(tctx *TraceContext, startTime time.Time, err error, depth int) {
	tctx.mtx.Lock()
	tctx.function = getFunc(depth)
	tctx.metadata["Latency"] = strconv.FormatInt(time.Since(startTime).Nanoseconds()/1000000, 10)
	if err != nil {
		Logger.Print(timePrefix() + " Level=\"" + errorLog + "\" Msg=\"EndTrace\" Err=\"" + err.Error() + "\"" + convertTags(tctx))
	} else {
		Logger.Print(timePrefix() + " Level=\"" + infoLog + "\" Msg=\"EndTrace\"" + convertTags(tctx))
	}
	tctx.mtx.Unlock()
}
