package logging

import (
	"encoding/json"
	"fmt"
	"strings"
)

// LogLevel is a representation of the logging level for a logger.
//go:generate stringer -type=LogLevel
type LogLevel int

// Definition of LogLevels for a logger.
const (
	_ LogLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
)

// GetLogLevelsForString will get the appropriate loglevels for a string
// log level representation.
func GetLogLevelsForString(logLevel string) ([]LogLevel, error) {
	var logLevels []LogLevel
	switch strings.ToLower(logLevel) {
	case "debug":
		logLevels = append(logLevels, DEBUG)
		fallthrough
	case "info":
		logLevels = append(logLevels, INFO)
		fallthrough
	case "warn":
		logLevels = append(logLevels, WARN)
		fallthrough
	case "error":
		logLevels = append(logLevels, ERROR)
	default:
		return nil, fmt.Errorf(
			"couldn't find loglevel for string '%s'",
			logLevel,
		)
	}

	return logLevels, nil
}

// MarshalJSON correctly formats log into json format.
func (ll LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(ll.String())
}
