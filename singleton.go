package slogging

import (
	"os"
	"sync"
)

var loggersRWMutex = new(sync.RWMutex)
var allLoggers = make(map[string]Logger)

var defaultLoggerName = "default"

// GetDefaultLogger gets the default logger.
func GetDefaultLogger() Logger {
	loggersRWMutex.RLock()
	if logger, ok := allLoggers[defaultLoggerName]; ok {
		return logger
	}
	loggersRWMutex.RUnlock()

	logLevel := os.Getenv("SLOGGING_DEFAULT_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}

	logLevels, err := GetLogLevelsForString(logLevel)
	if err != nil {
		panic(err)
	}

	logger := GetJSONLogger(
		Stdout,
		logLevels,
	)

	loggersRWMutex.Lock()
	allLoggers[defaultLoggerName] = &logger
	loggersRWMutex.Unlock()

	return &logger
}

// SetDefaultLogger will use the provided identifier as the default
// logger for future log calls.
func SetDefaultLogger(identifier string, logger Logger) {
	loggersRWMutex.Lock()
	defer loggersRWMutex.Unlock()
	defaultLoggerName = identifier
	allLoggers[identifier] = logger
}

// SetDefaultLoggerName will use the provided identifier as the default
// logger for future log calls.
func SetDefaultLoggerName(identifier string) {
	loggersRWMutex.Lock()
	defer loggersRWMutex.Unlock()
	defaultLoggerName = identifier
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
