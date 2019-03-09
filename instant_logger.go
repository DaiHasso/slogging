package logging

import (
    "encoding/json"
    "fmt"
    "os"
    "io"
    "log"
    "regexp"
    "strings"
    "time"

    "github.com/pkg/errors"
)

var standardReplaceRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)

type InstantLogger struct {
    writerLoggers map[io.Writer]*log.Logger
    logLevelsEnabled map[LogLevel]bool
    format LogFormat

    extraGenerators []ExtrasGenerator
}

func padStringRight(str, pad string, count int) string {
    result := str
    for i := 0; i < count; i++ {
        result += pad
    }

    return result
}

func formatStandardBodyLong(bodyMap map[string]interface{}) []byte {
    var (
        firstParts [2]string
        parts []string
        lastPart string
    )

    for key, value := range bodyMap {
        valueString := fmt.Sprint(value)

        // TODO: Maybe this logic goes somewhere else?
        if key == "message" {
            lastPart = valueString
        } else if key == "log_level" {
            firstParts[1] = valueString
        } else if key == "timestamp" {
            firstParts[0] = valueString
        } else {
            partString := fmt.Sprintf(
                `%s="%s"`, key, valueString,
            )
            parts = append(parts, partString)
        }
    }

    parts = append(firstParts[:], append(parts, lastPart)...)
    logLine := strings.Join(parts, " ")
    return []byte(logLine)
}

func formatStandardBody(bodyMap map[string]interface{}) []byte {
    var (
        keys,
        values []string
    )

    var (
        firstKeys,
        firstValues [3]string
    )

    for key, value := range bodyMap {
        valueString := fmt.Sprint(value)

        keyLen := len(key)
        valueLen := len(valueString)

        finalKey := key
        finalValue := valueString
        if keyLen < valueLen {
            finalKey = padStringRight(key, " ", valueLen - keyLen)
        } else if keyLen > valueLen {
            finalValue = padStringRight(valueString, " ", keyLen - valueLen)
        }

        // TODO: Maybe this logic goes somewhere else?
        if key == "message" {
            firstKeys[2] = finalKey
            firstValues[2] = finalValue
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

    keys = append(firstKeys[:], keys...)
    values = append(firstValues[:], values...)

    headerString := strings.Join(keys, " | ")
    valuesString := strings.Join(values, " | ")

    return []byte(headerString + "\n" + valuesString)
}

func (self InstantLogger) applyExtras(extras []Extras) map[string]interface{} {
    allExtras := make(map[string]interface{})
    for _, extra := range extras {
        for key, value := range extra {
            allExtras[key] = value
        }
    }

    return allExtras
}

func (self InstantLogger) formatMessage(
    logBody map[string]interface{},
) []byte {
    switch self.format {
    case JSON:
        bodyBytes, err := json.Marshal(logBody)
        if err == nil {
            return bodyBytes
        }
    case Standard:
        return formatStandardBody(logBody)
    case StandardLong:
        return formatStandardBodyLong(logBody)
    }

    // NOTE: Fallback on straight message if we have a problem with marshaling.
    // This is very unlikely to happen.
    return []byte(fmt.Sprintf(
        "Error while marshalling log with message '%s'.",
        logBody["message"],
    ))
}

func (self InstantLogger) addDefaults(
    logLevel LogLevel, message string, logBody map[string]interface{},
) map[string]interface{} {
    logBody["message"] = message
    logBody["log_level"] = logLevel
    logBody["timestamp"] = timestamp{time.Now()}

    return logBody
}

func (self InstantLogger) levelEnabled(level LogLevel) bool {
    _, ok := self.logLevelsEnabled[level]
    return ok
}

func (self InstantLogger) applyInstanceExtras() ([]Extras, error) {
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

func (self InstantLogger) applyGlobalExtras() ([]Extras, error) {
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

func (self InstantLogger) logToLevel(
    level LogLevel, message string, extras []Extras,
) {
    allExtras := extras
    extraGenerators, err := self.applyInstanceExtras()
    allExtras = append(allExtras, extraGenerators...)
    if err != nil {
        self.Exception(err, "Error while running logger instance extras.")
    }

    globalExtras, err := self.applyGlobalExtras()
    allExtras = append(allExtras, globalExtras...)
    if err != nil {
        self.Exception(err, "Error while running global logger extras.")
    }

    extrasMap := self.applyExtras(allExtras)
    extrasMap = self.addDefaults(level, message, extrasMap)

    logBody := self.formatMessage(extrasMap)

    self.Log(level, logBody)
}

func (self InstantLogger) Log(level LogLevel, messageBytes []byte) {
    if !self.levelEnabled(level) {
        return
    }

    for _, logger := range self.writerLoggers {
        logger.Printf("%s", messageBytes)
    }
}

func (self InstantLogger) Debug(message string, extras ...Extras) {
    self.logToLevel(DEBUG, message, extras)
}

func (self InstantLogger) Info(message string, extras ...Extras) {
    self.logToLevel(INFO, message, extras)
}

func (self InstantLogger) Warn(message string, extras ...Extras) {
    self.logToLevel(WARN, message, extras)
}

func (self InstantLogger) Error(message string, extras ...Extras) {
    self.logToLevel(ERROR, message, extras)
}

func (self InstantLogger) Exception(
    err error, message string, extras ...Extras,
) {
    extrasWithErr := append(extras, Extras{
        "error":  fmt.Sprintf("%+v", errors.WithStack(err)),
    })

    self.logToLevel(ERROR, message, extrasWithErr)
}

func (self *InstantLogger) AddDefaultExtras(
    extras ExtrasGenerator, otherExtras ...ExtrasGenerator,
) {
    allExtras := append([]ExtrasGenerator{extras}, otherExtras...)
    self.extraGenerators = append(self.extraGenerators, allExtras...)
}

func (self *InstantLogger) SetDefaultExtras(
    extras ExtrasGenerator, otherExtras ...ExtrasGenerator,
) {
    allExtras := append([]ExtrasGenerator{extras}, otherExtras...)
    self.extraGenerators = allExtras
}

func (self *InstantLogger) SetFormat(logFormat LogFormat) {
    self.format = logFormat
}

func (self *InstantLogger) SetWriters(w io.Writer, otherWs ...io.Writer) {
    newWriters := make(map[io.Writer]*log.Logger)
    for _, writer := range append([]io.Writer{w}, otherWs...) {
        newWriters[writer] = log.New(writer, "", 0)
    }

    self.writerLoggers = newWriters
}

func (self *InstantLogger) AddWriters(w io.Writer, otherWs ...io.Writer) {
    for _, writer := range append([]io.Writer{w}, otherWs...) {
        self.writerLoggers[writer] = log.New(writer, "", 0)
    }
}

func (self *InstantLogger) RemoveWriter(w io.Writer) {
    if _, ok := self.writerLoggers[w]; ok {
        delete(self.writerLoggers, w)
    }
}

// NewLogger creates a new logger basing it off the default/root logger.
func NewLogger(
    identifier string, options ...LoggerOption,
) (*InstantLogger, error) {
    loggerConfig := newLoggerConfig()
    for i, opt := range(options) {
        err := opt(loggerConfig)
        if err != nil {
            return nil, errors.Wrapf(
                err, "Error while processing option #%d", i,
            )
        }
    }

    writerLoggers := loggerConfig.writerLoggers
    if len(writerLoggers) == 0 {
        writerLoggers[os.Stdout] = log.New(os.Stdout, "", 0)
    }

    logsEnabled := loggerConfig.logsEnabled
    if len(logsEnabled) == 0 {
        logsEnabled, _ = logsEnabledFromLevel(INFO)
    }

    format := loggerConfig.logFormat
    if format == UnsetFormat {
        format = Standard
    }

    instantLogger := &InstantLogger{
        writerLoggers: writerLoggers,
        logLevelsEnabled: logsEnabled,
        format: format,
        extraGenerators: loggerConfig.extraGenerators,
    }

    return instantLogger, nil
}
