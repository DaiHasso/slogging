package logging

import (
    "io"
    "log"

    "github.com/pkg/errors"
)

type loggerConfig struct{
    writerLoggers map[io.Writer]*log.Logger
    logsEnabled map[LogLevel]bool
    logFormat LogFormat
    extraGenerators []ExtrasGenerator
}

func newLoggerConfig() *loggerConfig {
    return &loggerConfig{
        writerLoggers: make(map[io.Writer]*log.Logger),
        logsEnabled: make(map[LogLevel]bool),
        logFormat: UnsetFormat,
        extraGenerators: make([]ExtrasGenerator, 0),
    }
}

// LoggerOption is an option used when creating a new Logger.
type LoggerOption func(*loggerConfig) error

// WithLogLevel sets this Logger's log level to the provided LogLevel.
func WithLogLevel(logLevel LogLevel) LoggerOption {
    return func(loggerConfig *loggerConfig) error {
        logsEnabled, err := logsEnabledFromLevel(logLevel)
        if err != nil {
            return errors.Wrapf(
                err,
                "Error setting log level from provided level string '%s'",
                logLevel,
            )
        }

        loggerConfig.logsEnabled = logsEnabled

        return nil
    }
}

// WithFormat sets the new Logger's format to the provided format.
func WithFormat(logFormat LogFormat) LoggerOption {
    return func(loggerConfig *loggerConfig) error {
        loggerConfig.logFormat = logFormat

        return nil
    }
}

// WithLogWriters sets the provided Logger's writers to the provided writers.
func WithLogWriters(primary io.Writer, others ... io.Writer) LoggerOption {
    return func(loggerConfig *loggerConfig) error {
        for _, writer := range append([]io.Writer{primary}, others...) {
            loggerConfig.writerLoggers[writer] = log.New(writer, "", 0)
        }

        return nil
    }
}

// WithDefaultExtras provides one or many Extras that will be logged for every
// log statement for this Logger.
func WithDefaultExtras(
    extraParam ExtrasGenerator, extraParams ...ExtrasGenerator,
) LoggerOption {
    allExtraParams := append([]ExtrasGenerator{extraParam}, extraParams...)
    return func(loggerConfig *loggerConfig) error {
        loggerConfig.extraGenerators = append(
            loggerConfig.extraGenerators, allExtraParams...,
        )

        return nil
    }
}
