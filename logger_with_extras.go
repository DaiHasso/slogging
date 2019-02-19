package logging

import (
    "github.com/pkg/errors"
)

// LoggerWithExtras wraps a logger and applies extras before starting a log
// chain.
type LoggerWithExtras struct {
    Logger

    extras []ExtraParameter
}

func (self LoggerWithExtras) applyExtras(logInstance LogInstance) LogInstance {
    for i, extra := range self.extras {
        err := extra(logInstance)
        if err != nil {
            panic(errors.Wrapf(err, "Error while applying extra #%d", i))
        }
    }

    return logInstance
}

func (self LoggerWithExtras) Debug(message string) LogInstance {
    logInstance := self.Logger.Debug(message)
    return self.applyExtras(logInstance)
}

func (self LoggerWithExtras) Warn(message string) LogInstance {
    logInstance := self.Logger.Warn(message)
    return self.applyExtras(logInstance)
}

func (self LoggerWithExtras) Error(message string) LogInstance {
    logInstance := self.Logger.Error(message)
    return self.applyExtras(logInstance)
}

func (self LoggerWithExtras) Info(message string) LogInstance {
    logInstance := self.Logger.Info(message)
    return self.applyExtras(logInstance)
}

// NewLoggerWithExtras wraps an existing logger adding extras to the logger.
func NewLoggerWithExtras(
    existing Logger, extras ...ExtraParameter,
) LoggerWithExtras {
    return LoggerWithExtras{
        Logger: existing,
        extras: extras,
    }
}
