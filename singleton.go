package slogging

import "os"

var allLoggers = make(map[string]Logger)

var defaultLoggerName = "default"

// GetDefaultLogger gets the default logger.
func GetDefaultLogger() Logger {
	if logger, ok := allLoggers[defaultLoggerName]; ok {
		return logger
	}

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

	return &logger
}

// SetDefaultLogger will use the provided identifier as the default
// logger for future log calls.
func SetDefaultLogger(identifier string, logger Logger) {
	defaultLoggerName = identifier
	allLoggers[identifier] = logger
}

// SetDefaultLoggerName will use the provided identifier as the default
// logger for future log calls.
func SetDefaultLoggerName(identifier string) {
	defaultLoggerName = identifier
}

// Debug uses the default logger to log to debug level.
func Debug(message string) LogInstance {
	logger := getDefaultLogger()

	return logger.Debug(message)
}

// Warn uses the default logger to log to warn level.
func Warn(message string) LogInstance {
	logger := getDefaultLogger()

	return logger.Warn(message)
}

// Error uses the default logger to log to error level.
func Error(message string) LogInstance {
	logger := getDefaultLogger()

	return logger.Error(message)
}

// Info uses the default logger to log to info level.
func Info(message string) LogInstance {
	logger := getDefaultLogger()

	return logger.Info(message)
}
