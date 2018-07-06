package slogging

import (
	"fmt"
	"runtime"
)

func traceError() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[1])
	file, line := f.FileLine(pc[1])
	return fmt.Sprintf("%s:%d %s", file, line, f.Name())
}
