package os_utils

import "os"

func Exit(code int, isIgnore bool) {
	if !isIgnore {
		os.Exit(code)
	}
}
