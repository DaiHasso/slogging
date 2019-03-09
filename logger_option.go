package logging

import (
    "io"
    "os"
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

type LoggerOption func(*loggerConfig) error

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

func WithFormat(logFormat LogFormat) LoggerOption {
    return func(loggerConfig *loggerConfig) error {
        loggerConfig.logFormat = logFormat

        return nil
    }
}

func WithTarget(
    logTarget LogTarget, otherLogTargets ...LogTarget,
) LoggerOption {
    return func(loggerConfig *loggerConfig) error {
        targets := append([]LogTarget{logTarget}, otherLogTargets...)
        for _, target := range targets {
            if target == Stdout {
                loggerConfig.writerLoggers[os.Stdout] = log.New(
                    os.Stdout, "", 0,
                )
            }
        }

        return nil
    }
}

func WithLogWriters(primary io.Writer, others ... io.Writer) LoggerOption {
    return func(loggerConfig *loggerConfig) error {
        for _, writer := range append([]io.Writer{primary}, others...) {
            loggerConfig.writerLoggers[writer] = log.New(writer, "", 0)
        }

        return nil
    }
}

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
