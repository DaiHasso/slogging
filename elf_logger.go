package logging

import (
	"log"
	"os"
)

// ELFLogger is an instance of a logger in Extended Log Format (ELF).
// More info here: https://www.w3.org/TR/WD-logfile.html
type ELFLogger struct {
	logTarget   LogTarget
	logger      *log.Logger
	logsEnabled map[LogLevel]bool
}

// Debug starts a new log at debug level.
func (el *ELFLogger) Debug(message string) LogInstance {
	return getELFLog(el, DEBUG, message)
}

// Warn starts a new log at warn level.
func (el *ELFLogger) Warn(message string) LogInstance {
	return getELFLog(el, WARN, message)
}

// Error starts a new log at error level.
func (el *ELFLogger) Error(message string) LogInstance {
	return getELFLog(el, ERROR, message)
}

// Info starts a new log at info level.
func (el *ELFLogger) Info(message string) LogInstance {
	return getELFLog(el, INFO, message)
}

// Log logs a specific message at a specified log level.
func (el *ELFLogger) Log(
	logLevel LogLevel,
	outputBytes []byte,
) {
	if el.logsEnabled[logLevel] {
		el.logger.Printf(
			"%s\n",
			outputBytes,
		)
	}
}

// SetInternalLogger sets the logger used internally, mostly for testing.
func (el *ELFLogger) SetInternalLogger(logger *log.Logger) {
	el.logger = logger
}

// GetInternalLogger will return the logger used internally, mostly for
// testing.
func (el *ELFLogger) GetInternalLogger() *log.Logger {
	return el.logger
}

// GetPseudoWriter get a PseudoWriter for a ELFLogger.
func (el *ELFLogger) GetPseudoWriter(logLevel LogLevel) PseudoWriter {
	return PseudoWriter{
		logger:   el,
		logLevel: logLevel,
	}
}

// SetPretty has no effect on an ELFLogger.
func (el *ELFLogger) SetPretty(pretty bool) {}

// GetStdLogger returns a std library logger that uses
// a PseudoWriter.
func (el *ELFLogger) GetStdLogger(logLevel LogLevel) *log.Logger {
	return log.New(
		el.GetPseudoWriter(logLevel),
		"",
		0,
	)
}

// SetLogLevel will set the log level for this logger to the provided level.
func (el *ELFLogger) SetLogLevel(logLevel string) error {
	logLevels, err := GetLogLevelsForString(logLevel)
	if err != nil {
		return err
	}

	logsEnabledMap := make(map[LogLevel]bool)
	for _, level := range logLevels {
		logsEnabledMap[level] = true
	}

	el.logsEnabled = logsEnabledMap

	return nil
}

// GetELFLogger will instantiate a new ELF logger writing to the specified log
// target and with the specified log levels enabled.
func GetELFLogger(
	logTarget LogTarget,
	logsEnabled []LogLevel,
) *ELFLogger {
	var logger *log.Logger

	if logTarget == Stdout {
		logger = log.New(os.Stdout, "", 0)
	}

	logsEnabledMap := make(map[LogLevel]bool)
	for _, level := range logsEnabled {
		logsEnabledMap[level] = true
	}

	newLogger := &ELFLogger{
		logger:      logger,
		logTarget:   logTarget,
		logsEnabled: logsEnabledMap,
	}

	return newLogger
}
