package slogging

import (
	"fmt"
	"log"
)

// Logger is the base interface for all loggers.
type Logger interface {
	Debug(string) LogInstance
	Warn(string) LogInstance
	Error(string) LogInstance
	Info(string) LogInstance
	Log(LogLevel, []byte)
	SetPretty(bool)
	SetInternalLogger(*log.Logger)
	GetInternalLogger() *log.Logger
	GetPseudoWriter(LogLevel) PseudoWriter
	GetStdLogger(LogLevel) *log.Logger
}

// GetNewLogger will get a new logger with the specified format,
// target and enabled logs then add it to the global log list.
func GetNewLogger(
	identifier string,
	logFormat LogFormat,
	logTarget LogTarget,
	logsEnabled []LogLevel,
) Logger {
	switch logFormat {
	case JSON:
		logger := GetJSONLogger(logTarget, logsEnabled)
		allLoggers[identifier] = &logger
		return &logger
	case ELF, Standard:
		logger := GetELFLogger(logTarget, logsEnabled)
		allLoggers[identifier] = &logger
		return &logger
	}

	panic(fmt.Errorf(
		"no implementation for LogFormat: '%s'",
		string(logFormat),
	))
}
