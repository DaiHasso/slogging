package logging

import (
    "fmt"
    "strings"
)

// LogLevel is a representation of the logging level for a logger.
//go:generate stringer -type=LogLevel
type LogLevel string

// Definition of LogLevels for a logger.
const (
    UnsetLogLevel LogLevel = ""
    ERROR LogLevel = "ERROR"
    WARN LogLevel = "WARN"
    INFO LogLevel = "INFO"
    DEBUG LogLevel = "DEBUG"
)

// GetLogLevelsForString will get the appropriate loglevels for a string
// log level representation.
func GetLogLevelsForString(logLevel string) (map[LogLevel]bool, error) {
    return logsEnabledFromLevel(LogLevel(strings.ToUpper(logLevel)))
}

func logsEnabledFromLevel(logLevel LogLevel) (map[LogLevel]bool, error) {
    logLevels := make(map[LogLevel]bool)
    switch logLevel {
    case DEBUG:
        logLevels[DEBUG] = true
        fallthrough
    case INFO:
        logLevels[INFO] = true
        fallthrough
    case WARN:
        logLevels[WARN] = true
        fallthrough
    case ERROR:
        logLevels[ERROR] = true
    default:
        return nil, fmt.Errorf(
            "Incorrect LogLevel: '%s'",
            logLevel,
        )
    }

    return logLevels, nil
}

func LogLevelFromString(logLevel string) LogLevel {
    return LogLevel(strings.ToUpper(logLevel))
}
