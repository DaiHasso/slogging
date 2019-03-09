package logging

import (
    "github.com/pkg/errors"
)

// LoggerWithExtras wraps a logger and applies extras before starting a log
// chain.
type LoggerWithExtras struct {
    ChainLogger

    extras []ExtrasGenerator
}

func (self LoggerWithExtras) applyExtras(logInstance LogInstance) LogInstance {
    err := applyExtrasToLogInstance(logInstance, self.extras)
    if err != nil {
        panic(errors.Wrap(err, "Error while applying extra"))
    }
    return logInstance
}

func (self LoggerWithExtras) Debug(message string) LogInstance {
    logInstance := self.ChainLogger.Debug(message)
    return self.applyExtras(logInstance)
}

func (self LoggerWithExtras) Warn(message string) LogInstance {
    logInstance := self.ChainLogger.Warn(message)
    return self.applyExtras(logInstance)
}

func (self LoggerWithExtras) Error(message string) LogInstance {
    logInstance := self.ChainLogger.Error(message)
    return self.applyExtras(logInstance)
}

func (self LoggerWithExtras) Info(message string) LogInstance {
    logInstance := self.ChainLogger.Info(message)
    return self.applyExtras(logInstance)
}

// NewLoggerWithExtras wraps an existing logger adding extras to the logger.
func NewLoggerWithExtras(
    existing ChainLogger, extras ...ExtrasGenerator,
) LoggerWithExtras {
    return LoggerWithExtras{
        ChainLogger: existing,
        extras: extras,
    }
}
