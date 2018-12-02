package slogging

import (
	"log"
	"os"
)

// JSONLogger is a logger that logs output in JSON.
type JSONLogger struct {
	logTarget   LogTarget
	logger      *log.Logger
	logsEnabled map[LogLevel]bool
	pretty      bool
}

// Debug returns a jsonLog which is used to build the output
// for the log.
func (jl *JSONLogger) Debug(message string) LogInstance {
	return getJSONLog(jl, DEBUG, message, jl.pretty)
}

// Warn returns a jsonLog which is used to build the output
// for the log.
func (jl *JSONLogger) Warn(message string) LogInstance {
	return getJSONLog(jl, WARN, message, jl.pretty)
}

// Error returns a jsonLog which is used to build the output
// for the log.
func (jl *JSONLogger) Error(message string) LogInstance {
	jsonLog := getJSONLog(jl, ERROR, message, jl.pretty).
		With("trace", traceError())
	return jsonLog
}

// Info returns a jsonLog which is used to build the output
// for the log.
func (jl *JSONLogger) Info(message string) LogInstance {
	return getJSONLog(jl, INFO, message, jl.pretty)
}

//Log takes a log level and an output string and logs the outputString.
func (jl *JSONLogger) Log(logLevel LogLevel, outputString []byte) {
	if jl.logsEnabled[logLevel] {
		jl.logger.Printf(
			"%s\n",
			outputString,
		)
	}
}

// SetInternalLogger sets the logger used internally, mostly for testing.
func (jl *JSONLogger) SetInternalLogger(logger *log.Logger) {
	jl.logger = logger
}

// GetInternalLogger will return the logger used internally, mostly for
// testing.
func (jl *JSONLogger) GetInternalLogger() *log.Logger {
	return jl.logger
}

// GetPseudoWriter get a PseudoWriter for a JSONLogger.
func (jl *JSONLogger) GetPseudoWriter(logLevel LogLevel) PseudoWriter {
	return PseudoWriter{
		logger:   jl,
		logLevel: logLevel,
	}
}

// SetPretty will set the pretty-print setting.
func (jl *JSONLogger) SetPretty(pretty bool) {
	jl.pretty = pretty
}

// GetStdLogger returns a std library logger that uses
// a PseudoWriter.
func (jl *JSONLogger) GetStdLogger(logLevel LogLevel) *log.Logger {
	return log.New(
		jl.GetPseudoWriter(logLevel),
		"",
		0,
	)
}

// SetLogLevel will set the log level for this logger to the provided level.
func (jl *JSONLogger) SetLogLevel(logLevel string) error {
	logLevels, err := GetLogLevelsForString(logLevel)
	if err != nil {
		return err
	}

	logsEnabledMap := make(map[LogLevel]bool)
	for _, level := range logLevels {
		logsEnabledMap[level] = true
	}

	jl.logsEnabled = logsEnabledMap

	return nil
}

// GetJSONLogger initializes and returns a new JSON logger.
func GetJSONLogger(
	logTarget LogTarget,
	logsEnabled []LogLevel,
) *JSONLogger {
	var logger *log.Logger

	if logTarget == Stdout {
		logger = log.New(os.Stdout, "", 0)
	}

	logsEnabledMap := make(map[LogLevel]bool)
	for _, level := range logsEnabled {
		logsEnabledMap[level] = true
	}

	newLogger := &JSONLogger{
		logger:      logger,
		logTarget:   logTarget,
		logsEnabled: logsEnabledMap,
		pretty:      false,
	}

	return newLogger
}
