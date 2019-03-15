package logging

import (
    "fmt"
    "io"
    "log"
    "os"
    "sync"

    "github.com/pkg/errors"
)

var (
    once sync.Once
    globalExtraGeneratorsMutex,
    loggersRWMutex,
    rootLoggerRWMutex *sync.RWMutex

    allLoggers map[string]*Logger

    globalExtraGenerators []ExtrasGenerator

    rootLoggerName string
    initialRootLoggerName = "root"

    SloggingEnvVarPrefix = "SLOGGING"
    prefixEnvVar = func(s string) string {
        return fmt.Sprintf("%s_%s", SloggingEnvVarPrefix, s)
    }
    RootLoggerLevelEnvVar = prefixEnvVar("ROOT_LOGGER_LEVEL")
    RootLoggerFormatEnvVar = prefixEnvVar("ROOT_LOGGER_FORMAT")
)

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

// GetLogger get an existing logger by its identifier.
func GetLogger(identifier string) *Logger {
    loggersRWMutex.RLock()
    defer loggersRWMutex.RUnlock()
    return allLoggers[identifier]
}

// GetRootLogger gets the root logger.
func GetRootLogger() *Logger {
    loggersRWMutex.RLock()
    rootLoggerRWMutex.RLock()
    defer loggersRWMutex.RUnlock()
    defer rootLoggerRWMutex.RUnlock()
    logger := allLoggers[rootLoggerName]

    return logger
}

// SetRootLogger will use the provided identifier as the default
// logger for future log calls.
func SetRootLogger(identifier string, logger *Logger) error {
    err := addLogger(identifier, logger)
    if err != nil {
        return errors.Wrap(err, "Error while trying to set root logger")
    }

    rootLoggerRWMutex.Lock()
    defer rootLoggerRWMutex.Unlock()
    rootLoggerName = identifier

    return nil
}

// SetRootLoggerExisting will use the provided identifier as the default
// logger for future log calls.
func SetRootLoggerExisting(identifier string) error {
    loggersRWMutex.Lock()
    defer loggersRWMutex.Unlock()
    if _, ok := allLoggers[identifier]; !ok {
        return errors.Errorf(
            "Logger with provided identifier '%s' doesn't exist.",
            identifier,
        )
    }

    rootLoggerRWMutex.Lock()
    defer rootLoggerRWMutex.Unlock()
    rootLoggerName = identifier

    return nil
}

// Debug uses the root logger to log to debug level.
func Debug(message string, extras ...Extras) {
    logger := GetRootLogger()

    loggersRWMutex.RLock()
    defer loggersRWMutex.RUnlock()

    logger.Debug(message, extras...)
}

// Warn uses the root logger to log to warn level.
func Warn(message string, extras ...Extras) {
    logger := GetRootLogger()

    loggersRWMutex.RLock()
    defer loggersRWMutex.RUnlock()

    logger.Warn(message, extras...)
}

// Error uses the root logger to log to error level.
func Error(message string, extras ...Extras) {
    logger := GetRootLogger()

    loggersRWMutex.RLock()
    defer loggersRWMutex.RUnlock()

    logger.Error(message, extras...)
}

// Info uses the root logger to log to info level.
func Info(message string, extras ...Extras) {
    logger := GetRootLogger()

    loggersRWMutex.RLock()
    defer loggersRWMutex.RUnlock()

    logger.Info(message, extras...)
}

// Exception uses the root logger to log an error at error level.
func Exception(err error, message string, extras ...Extras) {
    logger := GetRootLogger()

    loggersRWMutex.RLock()
    defer loggersRWMutex.RUnlock()

    logger.Exception(err, message, extras...)
}

func addLogger(identifier string, logger *Logger) error {
    if identifier == "" {
        return errors.New("Identifier cannot be empty")
    }
    loggersRWMutex.Lock()
    defer loggersRWMutex.Unlock()
    if existing, ok := allLoggers[identifier]; ok {
        if existing != logger {
            return errors.New("Identifier already used for a different logger")
        }
        // No action needed if it's already in the map.
    } else {
        allLoggers[identifier] = logger
    }

    return nil
}

func removeLogger(identifier string) {
    loggersRWMutex.Lock()
    defer loggersRWMutex.Unlock()
    delete(allLoggers, identifier)
}

func identifierExists(identifier string) bool {
    loggersRWMutex.RLock()
    defer loggersRWMutex.RUnlock()
    _, ok := allLoggers[identifier]
    return ok
}

func init() {
    once.Do(func() {
        logLevel := INFO
        if logLevelString, ok := os.LookupEnv(RootLoggerLevelEnvVar); ok {
            logLevel = LogLevelFromString(logLevelString)
        }

        tempLogger := log.New(os.Stdout, "Slogging: ", 0)
        tempDebugLog := func(log string) {
            if logLevel == DEBUG {
                tempLogger.Println(log)
            }
        }

        // TODO: Be more clever about this?
        if logLevel == DEBUG {
            tempDebugLog("Init started.")
        }


        format := JSON
        if formatString, ok := os.LookupEnv(RootLoggerFormatEnvVar); ok {
            envFormat := FormatFromString(formatString)
            if envFormat == UnknownFormat {
                tempDebugLog(fmt.Sprintf(
                    "Found value for '%s' but it's value '%s' was not an " +
                        "understood format.",
                    RootLoggerFormatEnvVar,
                    formatString,
                ))
            } else {
                format = envFormat
            }
        }

        loggersRWMutex = new(sync.RWMutex)
        globalExtraGeneratorsMutex = new(sync.RWMutex)
        rootLoggerRWMutex = new(sync.RWMutex)

        rootLoggerRWMutex.Lock()
        rootLoggerName = initialRootLoggerName

        logLevelsEnabled, err := logsEnabledFromLevel(logLevel)
        if err != nil {
            rootLoggerRWMutex.Unlock()
            panic(err)
        }
        rootLogger := &Logger{
            identifier: rootLoggerName,
            format: format,
            writerLoggers: map[io.Writer]*log.Logger{
                os.Stdout: log.New(os.Stdout, "", 0),
            },
            logLevelsEnabled: logLevelsEnabled,
        }

        loggersRWMutex.Lock()
        allLoggers = map[string]*Logger{
            rootLoggerName: rootLogger,
        }
        rootLoggerRWMutex.Unlock()
        loggersRWMutex.Unlock()

        rootLogger.Debug("Slogging init end.")
    })
}
