package logging

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "regexp"
    "strings"
    "time"

    "github.com/pkg/errors"
)

var keySantizationRegexp = regexp.MustCompile(`[\n\r\s]`)

// Logger is a logger instance that provides a unified interface for logging
// data.
type Logger struct {
    identifier string
    writerLoggers map[io.Writer]*log.Logger
    logLevelsEnabled map[LogLevel]bool
    format LogFormat
    extraGenerators []ExtrasGenerator
}

// Used for standard formats so you don't get super weird logs. All bets are
// off for JSON though.
func sanitizeKey(key string) string {
    return keySantizationRegexp.ReplaceAllString(strings.TrimSpace(key), "_")
}

func padStringRight(str, pad string, count int) string {
    result := str
    for i := 0; i < count; i++ {
        result += pad
    }

    return result
}

func stringify(in interface{}) (string, bool) {
    switch val := in.(type) {
    case ([]byte):
        if StringifyByteArrays {
            return string(val), true
        }
    case error:
        return fmt.Sprintf("%#+v", val), true
    }

    return "", false
}

func formatStandardBody(bodyMap map[string]interface{}) []byte {
    var (
        firstParts [2]string
        parts []string
        lastPart string
    )

    for key, value := range bodyMap {
        var valueString string
        if stringVal, ok := stringify(value); ok {
            valueString = stringVal
        } else {
            valueString = fmt.Sprint(value)
        }

        // TODO: Maybe this logic goes somewhere else?
        if key == "message" {
            lastPart = valueString
        } else if key == "log_level" {
            firstParts[1] = valueString
        } else if key == "timestamp" {
            firstParts[0] = valueString
        } else {
            partString := fmt.Sprintf(
                `%s="%s"`, sanitizeKey(key), valueString,
            )
            parts = append(parts, partString)
        }
    }

    parts = append(firstParts[:], append(parts, lastPart)...)
    logLine := strings.Join(parts, " ")
    return []byte(logLine)
}

func formatJsonBody(bodyMap map[string]interface{}) ([]byte, error) {
    for k, v := range bodyMap {
        if stringVal, ok := stringify(v); ok {
            bodyMap[k] = stringVal
        }
    }

    bodyBytes, err := json.Marshal(bodyMap)
    if err != nil {
        return nil, err
    }
    return bodyBytes, nil
}

func formatStandardBodyExtended(bodyMap map[string]interface{}) []byte {
    var (
        keys,
        values []string
        firstKeys,
        firstValues [2]string
        lastKeys,
        lastValues [1]string
    )

    for key, value := range bodyMap {
        valueString := fmt.Sprint(value)

        keyLen := len(key)
        valueLen := len(valueString)

        finalKey := sanitizeKey(key)
        finalValue := valueString
        if keyLen < valueLen {
            finalKey = padStringRight(key, " ", valueLen - keyLen)
        } else if keyLen > valueLen {
            finalValue = padStringRight(valueString, " ", keyLen - valueLen)
        }

        // TODO: Maybe this logic goes somewhere else?
        if key == "message" {
            lastKeys[0] = finalKey
            lastValues[0] = finalValue
        } else if key == "log_level" {
            firstKeys[1] = finalKey
            firstValues[1] = finalValue
        } else if key == "timestamp" {
            firstKeys[0] = finalKey
            firstValues[0] = finalValue
        } else {
            keys = append(keys, finalKey)
            values = append(values, finalValue)
        }
    }

    keys = append(append(firstKeys[:], keys...), lastKeys[:]...)
    values = append(append(firstValues[:], values...), lastValues[:]...)

    headerString := strings.Join(keys, " | ")
    valuesString := strings.Join(values, " | ")

    return []byte(headerString + "\n" + valuesString)
}

func (self Logger) clone() *Logger {
    newWriterLoggers := make(map[io.Writer]*log.Logger)
    for w, l := range self.writerLoggers {
        newWriterLoggers[w] = l
    }
    newLogsEnabled := make(map[LogLevel]bool)
    for logLevel, _ := range self.logLevelsEnabled {
        newLogsEnabled[logLevel] = true
    }
    return &Logger{
        identifier: "",
        writerLoggers: newWriterLoggers,
        logLevelsEnabled: newLogsEnabled,
        format: self.format,
        extraGenerators: self.extraGenerators,
    }
}


func (self Logger) applyExtras(extras []Extras) map[string]interface{} {
    allExtras := make(map[string]interface{})
    for _, extra := range extras {
        for key, value := range extra {
            allExtras[key] = value
        }
    }

    return allExtras
}

func (self Logger) formatMessage(
    logBody map[string]interface{},
) []byte {
    switch self.format {
    case JSON:
        bodyBytes, err := formatJsonBody(logBody)
        if err == nil {
            return bodyBytes
        }
    case Standard:
        return formatStandardBody(logBody)
    case StandardExtended:
        return formatStandardBodyExtended(logBody)
    }

    // NOTE: Fallback on straight message if we have a problem with marshaling.
    // This is very unlikely to happen.
    return []byte(fmt.Sprintf(
        "Error while marshalling log with message '%s'.",
        logBody["message"],
    ))
}

func (self Logger) addDefaults(
    logLevel LogLevel, message string, logBody map[string]interface{},
) map[string]interface{} {
    logBody["message"] = message
    logBody["log_level"] = logLevel
    logBody["timestamp"] = timestamp{time.Now()}

    return logBody
}

func (self Logger) levelEnabled(level LogLevel) bool {
    _, ok := self.logLevelsEnabled[level]
    return ok
}

func (self Logger) applyInstanceExtras() ([]Extras, error) {
    var allExtras []Extras
    for i, extraFunc := range self.extraGenerators {
        newExtras, err := extraFunc()
        if err != nil {
            return nil, errors.Wrapf(
                err, "Error while running extra #%d", i,
            )
        }
        allExtras = append(allExtras, newExtras)
    }

    return allExtras, nil
}

func (self Logger) applyGlobalExtras() ([]Extras, error) {
    var allExtras []Extras
    for i, extraFunc := range globalExtraGenerators {
        newExtras, err := extraFunc()
        if err != nil {
            return nil, errors.Wrapf(
                err, "Error while running extra #%d", i,
            )
        }
        allExtras = append(allExtras, newExtras)
    }

    return allExtras, nil
}

func (self Logger) internalException(err error, message string) {
    extras := Extras{
        "error":  fmt.Sprintf("%+v", errors.WithStack(err)),
    }

    body := self.addDefaults(ERROR, message, extras)

    logline := self.formatMessage(body)

    self.Log(ERROR, logline)
}

func (self Logger) logToLevel(
    level LogLevel, message string, extras []Extras,
) {
    if !self.levelEnabled(level) {
        // NOTE: It's definitely more efficient if we short-circuit here but is
        //       it proper?
        return
    }

    allExtras := extras
    extraGenerators, err := self.applyInstanceExtras()
    allExtras = append(allExtras, extraGenerators...)
    if err != nil {
        self.internalException(
            err, "Error while running logger instance extras.",
        )
    }

    globalExtras, err := self.applyGlobalExtras()
    allExtras = append(allExtras, globalExtras...)
    if err != nil {
        self.internalException(
            err, "Error while running global logger extras.",
        )
    }

    extrasMap := self.applyExtras(allExtras)
    extrasMap = self.addDefaults(level, message, extrasMap)

    logBody := self.formatMessage(extrasMap)

    self.Log(level, logBody)
}

// Log is the most basic log function. It logs the bytes directly if the
// loglevel is enabled. No aditional formating is done.
func (self Logger) Log(level LogLevel, messageBytes []byte) {
    if !self.levelEnabled(level) {
        return
    }

    for _, logger := range self.writerLoggers {
        logger.Printf("%s", messageBytes)
    }
}

// Debug logs according to this loggers formatter at the DEBUG level.
func (self Logger) Debug(message string, extras ...Extras) {
    self.logToLevel(DEBUG, message, extras)
}

// Info logs according to this loggers formatter at the INFO level.
func (self Logger) Info(message string, extras ...Extras) {
    self.logToLevel(INFO, message, extras)
}

// Warn logs according to this loggers formatter at the WARN level.
func (self Logger) Warn(message string, extras ...Extras) {
    self.logToLevel(WARN, message, extras)
}

// Error logs according to this loggers formatter at the ERROR level.
func (self Logger) Error(message string, extras ...Extras) {
    self.logToLevel(ERROR, message, extras)
}

// Exception logs an error's contents & stack at an error level.
func (self Logger) Exception(
    err error, message string, extras ...Extras,
) {
    extrasWithErr := append(extras, Extras{
        "error":  fmt.Sprintf("%+v", errors.WithStack(err)),
    })

    self.logToLevel(ERROR, message, extrasWithErr)
}

// AddDefaultExtras adds extra(s) which will be added for every log made with
// this logger.
func (self *Logger) AddDefaultExtras(
    extras ExtrasGenerator, otherExtras ...ExtrasGenerator,
) {
    allExtras := append([]ExtrasGenerator{extras}, otherExtras...)
    self.extraGenerators = append(self.extraGenerators, allExtras...)
}

// SetDefaultExtras sets (overriding) the extra(s) which will be added for
// every log made with this logger.
func (self *Logger) SetDefaultExtras(
    extras ExtrasGenerator, otherExtras ...ExtrasGenerator,
) {
    allExtras := append([]ExtrasGenerator{extras}, otherExtras...)
    self.extraGenerators = allExtras
}

// SetFormat changes the loggers format to the provided format.
func (self *Logger) SetFormat(logFormat LogFormat) {
    self.format = logFormat
}

// SetWriters sets the internal logger's writers to the provided writer(s).
func (self *Logger) SetWriters(w io.Writer, otherWs ...io.Writer) {
    newWriters := make(map[io.Writer]*log.Logger)
    for _, writer := range append([]io.Writer{w}, otherWs...) {
        newWriters[writer] = log.New(writer, "", 0)
    }

    self.writerLoggers = newWriters
}

// AddWriters adds writers provided to the existing writers if they don't
// already exist (duplicates will not be added multiple times).
func (self *Logger) AddWriters(w io.Writer, otherWs ...io.Writer) {
    for _, writer := range append([]io.Writer{w}, otherWs...) {
        self.writerLoggers[writer] = log.New(writer, "", 0)
    }
}

// RemoveWriter removes the provided writer if it is found.
func (self *Logger) RemoveWriter(w io.Writer) {
    if _, ok := self.writerLoggers[w]; ok {
        delete(self.writerLoggers, w)
    }
}

// SetLogLevel sets this logger to log at the provided LogLevel and below.
func (self *Logger) SetLogLevel(logLevel LogLevel) error {
    logsEnabled, err := logsEnabledFromLevel(logLevel)
    if err != nil {
        return errors.Wrap(err, "Error while setting log level")
    }

    self.logLevelsEnabled = logsEnabled

    return nil
}

// Identifier provides this Logger's identifier in the global Logger registry.
func (self Logger) Identifier() string {
    return self.identifier
}

// Close removes this logger from the global Logger registry and performs any
// cleanup tasks.
// It is not required to call this function when you're done with a logger but
// it is highly recommended to clear up memory and prevent accidental
// identifier clashing.
func (self Logger) Close() {
    removeLogger(self.identifier)
}

func newLogger(
    identifier string, baseLogger *Logger, options []LoggerOption,
) (*Logger, error) {
    if identifierExists(identifier) {
        return nil, errors.Errorf(
            "Can't create new logger with identifier '%s'; identifier " +
                "already exists",
            identifier,
        )
    }

    loggerConfig := newLoggerConfig()
    for i, opt := range(options) {
        err := opt(loggerConfig)
        if err != nil {
            return nil, errors.Wrapf(
                err, "Error while processing option #%d", i,
            )
        }
    }

    newLogger := baseLogger.clone()

    newLogger.identifier = identifier

    writerLoggers := loggerConfig.writerLoggers
    if len(writerLoggers) != 0 {
        newLogger.writerLoggers = writerLoggers
    }

    logsEnabled := loggerConfig.logsEnabled
    if len(logsEnabled) != 0 {
        newLogger.logLevelsEnabled = logsEnabled
    }

    format := loggerConfig.logFormat
    if format != UnsetFormat {
        newLogger.format = format
    }

    newLogger.extraGenerators = append(
        newLogger.extraGenerators, loggerConfig.extraGenerators...,
    )

    err := addLogger(identifier, newLogger)
    if err != nil {
        return nil, errors.Wrap(
            err, "Error while trying to add logger to global logger registry",
        )
    }

    return newLogger, nil
}

// NewLogger creates a new logger basing it off the default/root logger.
func NewLogger(
    identifier string, options ...LoggerOption,
) (*Logger, error) {
    return newLogger(identifier, GetRootLogger(), options)
}

// CloneLogger creates a new logger basing it off the provided logger.
func CloneLogger(
    identifier string, baseLogger *Logger, options ...LoggerOption,
) (*Logger, error) {
    return newLogger(identifier, baseLogger, options)
}
