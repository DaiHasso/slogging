package logging

import (
	"fmt"
	"strings"
)

// LogFormat is a representation of where a log should output to.
//go:generate stringer -type=LogFormat
type LogFormat int

// Definition of LogFormat for a logger.
const (
	_ LogFormat = iota
	JSON
	ELF
	Standard
)

// GetLogFormatForString will get the appropriate LogFormat for a string
// log format representation.
func GetLogFormatForString(logFormatString string) (LogFormat, error) {
	var logFormat LogFormat
	switch strings.ToLower(logFormatString) {
	case strings.ToLower(JSON.String()):
		logFormat = JSON
	case
		strings.ToLower(ELF.String()),
		strings.ToLower(Standard.String()):
		logFormat = ELF
	default:
		return -1, fmt.Errorf(
			"couldn't find LogFormat for string '%s'",
			logFormatString,
		)
	}

	return logFormat, nil
}
