package logging

import (
    "fmt"
    "log"
    "os"
    "sync"

    "github.com/pkg/errors"
)

var once sync.Once

var globalExtraGeneratorsMutex *sync.RWMutex
var loggersRWMutex *sync.RWMutex
var allLoggers map[string]ChainLogger

var initialDefaultLoggerName = "root"
var defaultLoggerName string
var globalExtraGenerators []ExtrasGenerator
var DefaultLogLevelEnvVar = "SLOGGING_DEFAULT_LOG_LEVEL"

func SetDefaultLoggerLogLevel(logLevel string) error {
    loggersRWMutex.Lock()
    defer loggersRWMutex.Unlock()
    logger := allLoggers[defaultLoggerName]
    return logger.SetLogLevel(logLevel)
}

// GetGlobalExtras returns the global extras.
func GetGlobalExtras() []ExtrasGenerator {
    globalExtraGeneratorsMutex.RLock()
    defer globalExtraGeneratorsMutex.RUnlock()
    return globalExtraGenerators
}

// SetGlobalExtras sets the global extras.
func SetGlobalExtras(extras ...ExtrasGenerator) {
    globalExtraGeneratorsMutex.Lock()
    defer globalExtraGeneratorsMutex.Unlock()
    globalExtraGenerators = extras
}

// AddGlobalExtras appends the provided extras to the global extras.
func AddGlobalExtras(extras ...ExtrasGenerator) {
    globalExtraGeneratorsMutex.Lock()
    defer globalExtraGeneratorsMutex.Unlock()
    globalExtraGenerators = append(globalExtraGenerators, extras...)
}

// GetDefaultLogger gets the default logger.
func GetDefaultLogger() ChainLogger {
    loggersRWMutex.RLock()
    defer loggersRWMutex.RUnlock()
    logger := allLoggers[defaultLoggerName]

    return logger
}

// SetDefaultLogger will use the provided identifier as the default
// logger for future log calls.
func SetDefaultLogger(identifier string, logger ChainLogger) {
    loggersRWMutex.Lock()
    defer loggersRWMutex.Unlock()
    defaultLoggerName = identifier
    allLoggers[identifier] = logger
}

// SetDefaultLoggerExisting will use the provided identifier as the default
// logger for future log calls.
func SetDefaultLoggerExisting(identifier string) error {
    loggersRWMutex.Lock()
    defer loggersRWMutex.Unlock()
    if _, ok := allLoggers[identifier]; !ok {
        return errors.Errorf(
            "Logger with provided identifier '%s' doesn't exist.",
            identifier,
        )
    }

    defaultLoggerName = identifier
    return nil
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

func loggerFromFormat(logFormat LogFormat) ChainLogger {
    switch logFormat {
    case JSON:
        logger := new(JSONLogger)
        return logger
    case ELF, Standard:
        logger := new(ELFLogger)
        return logger
    }

    panic(fmt.Errorf(
        "no implementation for LogFormat: '%s'",
        string(logFormat),
    ))
}

// NewChainLogger creates a new chain-style logger basing it off the
// default/root logger.
func NewChainLogger(
    identifier string, options ...LoggerOption,
) (ChainLogger, error) {
    loggerConfig := newLoggerConfig()
    for i, opt := range(options) {
        err := opt(loggerConfig)
        if err != nil {
            return nil, errors.Wrapf(
                err, "Error while processing option #%d", i,
            )
        }
    }

    logger := loggerFromFormat(loggerConfig.logFormat)

    GetDefaultLogger().CloneTo(logger)

    if len(loggerConfig.writerLoggers) != 0 {
        var first *log.Logger
        var rest []*log.Logger
        for _, logger := range loggerConfig.writerLoggers {
            if first == nil {
                first = logger
            } else {
                rest = append(rest, logger)
            }
        }
        logger.SetGoLoggers(first, rest...)
    }

    if loggerConfig.logsEnabled != nil {
        logger.SetLogsEnabled(loggerConfig.logsEnabled)
    }


    loggersRWMutex.Lock()
    allLoggers[identifier] = logger
    loggersRWMutex.Unlock()

    return logger, nil
}

// GetNewLogger will get a new logger with the specified format,
// target and enabled logs then add it to the global log list.
func GetNewLogger(
    identifier string,
    logFormat LogFormat,
    logTarget LogTarget,
    logsEnabled []LogLevel,
) ChainLogger {
    switch logFormat {
    case JSON:
        logger := GetJSONLogger(logTarget, logsEnabled)
        loggersRWMutex.Lock()
        allLoggers[identifier] = logger
        loggersRWMutex.Unlock()
        return logger
    case ELF, Standard:
        logger := GetELFLogger(logTarget, logsEnabled)
        allLoggers[identifier] = logger
        return logger
    }

    panic(fmt.Errorf(
        "no implementation for LogFormat: '%s'",
        string(logFormat),
    ))
}

func WithExtras(existing ChainLogger, extras ...ExtrasGenerator) ChainLogger {
    return NewLoggerWithExtras(existing, extras...)
}

func init() {
    once.Do(func() {
        logLevel := INFO
        if logLevelString, ok := os.LookupEnv(DefaultLogLevelEnvVar); ok {
            logLevel = LogLevelFromString(logLevelString)
        }

        // TODO: Be more clever about this.
        if logLevel == DEBUG {
            fmt.Println("Slogging init started.")
        }

        loggersRWMutex = new(sync.RWMutex)
        globalExtraGeneratorsMutex = new(sync.RWMutex)
        defaultLoggerName = initialDefaultLoggerName

        loggersRWMutex.Lock()
        newMap := make(map[string]ChainLogger)
        allLoggers = newMap

        logLevels, err := logsEnabledFromLevel(logLevel)
        if err != nil {
            logLevels, _ = logsEnabledFromLevel(INFO)
        }

        // TODO: Just use the new loggers by default
        var levelSlice []LogLevel
        for level, _ := range logLevels {
            levelSlice = append(levelSlice, level)
        }
        logger := GetJSONLogger(
            Stdout, levelSlice,
        )

        allLoggers[defaultLoggerName] = logger
        loggersRWMutex.Unlock()

        logger.Debug("Slogging init end.").Send()
    })
}
