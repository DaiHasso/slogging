package slogging

import (
	"fmt"
	"os"
	"sync"
	"strings"
)

var once sync.Once

var globalExtrasMutex *sync.RWMutex
var loggersRWMutex *sync.RWMutex
var allLoggers *map[string]Logger

var initialDefaultLoggerName = "default"
var defaultLoggerName *string
var globalExtras []ExtraParameter

func SetDefaultLoggerLogLevel(logLevel string) error {
	loggersRWMutex.Lock()
	defer loggersRWMutex.Unlock()
	logger := (*allLoggers)[*defaultLoggerName]
	return logger.SetLogLevel(logLevel)
}

// GetGlobalExtras returns the global extras.
func GetGlobalExtras() []ExtraParameter {
	globalExtrasMutex.RLock()
	defer globalExtrasMutex.RUnlock()
	return globalExtras
}

// SetGlobalExtras sets the global extras.
func SetGlobalExtras(extras ...ExtraParameter) {
	globalExtrasMutex.Lock()
	defer globalExtrasMutex.Unlock()
	globalExtras = extras
}

// AddGlobalExtras appends the provided extras to the global extras.
func AddGlobalExtras(extras ...ExtraParameter) {
	globalExtrasMutex.Lock()
	defer globalExtrasMutex.Unlock()
	globalExtras = append(globalExtras, extras...)
}

// GetDefaultLogger gets the default logger.
func GetDefaultLogger() Logger {
	loggersRWMutex.RLock()
	defer loggersRWMutex.RUnlock()
	logger := (*allLoggers)[*defaultLoggerName]

	return logger
}

// SetDefaultLogger will use the provided identifier as the default
// logger for future log calls.
func SetDefaultLogger(identifier string, logger Logger) {
	loggersRWMutex.Lock()
	defer loggersRWMutex.Unlock()
	defaultLoggerName = &identifier
	(*allLoggers)[identifier] = logger
}

// SetDefaultLoggerName will use the provided identifier as the default
// logger for future log calls.
func SetDefaultLoggerName(identifier string) {
	loggersRWMutex.Lock()
	defer loggersRWMutex.Unlock()
	defaultLoggerName = &identifier
}

// Debug uses the default logger to log to debug level.
func Debug(message string) LogInstance {
	logger := GetDefaultLogger()

	loggersRWMutex.RLock()
	defer loggersRWMutex.RUnlock()
	return logger.Debug(message)
}

// Warn uses the default logger to log to warn level.
func Warn(message string) LogInstance {
	logger := GetDefaultLogger()

	loggersRWMutex.RLock()
	defer loggersRWMutex.RUnlock()
	return logger.Warn(message)
}

// Error uses the default logger to log to error level.
func Error(message string) LogInstance {
	logger := GetDefaultLogger()

	loggersRWMutex.RLock()
	defer loggersRWMutex.RUnlock()
	return logger.Error(message)
}

// Info uses the default logger to log to info level.
func Info(message string) LogInstance {
	logger := GetDefaultLogger()

	loggersRWMutex.RLock()
	defer loggersRWMutex.RUnlock()
	return logger.Info(message)
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
		loggersRWMutex.Lock()
		(*allLoggers)[identifier] = logger
		loggersRWMutex.Unlock()
		return logger
	case ELF, Standard:
		logger := GetELFLogger(logTarget, logsEnabled)
		(*allLoggers)[identifier] = logger
		return logger
	}

	panic(fmt.Errorf(
		"no implementation for LogFormat: '%s'",
		string(logFormat),
	))
}

func init() {
	once.Do(func() {
		logLevel := strings.ToUpper(os.Getenv("SLOGGING_DEFAULT_LOG_LEVEL"))
		if logLevel == "" {
			logLevel = WARN.String()
		}

		// TODO: Be more clever about this.
		if logLevel == DEBUG.String() {
			fmt.Println("Slogging init started.")
		}

		loggersRWMutex = new(sync.RWMutex)
		globalExtrasMutex = new(sync.RWMutex)
		defaultLoggerName = &initialDefaultLoggerName

		loggersRWMutex.Lock()
		newMap := make(map[string]Logger)
		allLoggers = &newMap

		logLevels, err := GetLogLevelsForString(logLevel)
		if err != nil {
			panic(err)
		}

		logger := GetJSONLogger(
			Stdout,
			logLevels,
		)

		(*allLoggers)[*defaultLoggerName] = logger
		loggersRWMutex.Unlock()

		if logLevel == DEBUG.String() {
			fmt.Println("Slogging init end.")
		}
	})
}
