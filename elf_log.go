package slogging

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var version = 1.0
var versionHeader = fmt.Sprintf(
	"#Version: %s",
	strconv.FormatFloat(version, 'f', 1, 64),
)
var dateTemplate = "#Date: %s"
var fieldsTemplate = "#Fields: %s"
var formatTemplate = "%s\n%s\n%s\n%s\n"

// ELFLog is an instance of a log in Extended Log Format (ELF).
// More info here: https://www.w3.org/TR/WD-logfile.html
type ELFLog struct {
	elfLogger *ELFLogger
	contents  map[string]interface{}
	defaults  defaultValues
}

// SetTimestamp overrides timestamp, mostly used for testing.
func (el *ELFLog) SetTimestamp(newTime time.Time) LogInstance {
	el.defaults.Timestamp = timestamp{newTime}
	return el
}

// And performs the same functionality as @With.
func (el *ELFLog) And(key string, val interface{}) LogInstance {
	return el.With(key, val)
}

// With adds a key and value to a JSONLog.
func (el *ELFLog) With(key string, val interface{}) LogInstance {
	el.contents[strings.ToLower(key)] = val
	return el
}

// Send outputs the log entry to the logger and ensures proper
// formatting.
func (el *ELFLog) Send() {
	formattedOutput := el.formatForOutput()

	el.elfLogger.Log(
		el.defaults.Level,
		formattedOutput,
	)
}

func (el *ELFLog) formatForOutput() []byte {
	fieldsList := []string{"level"}
	fieldValuesList := []string{el.defaults.Level.String()}
	for key, value := range el.contents {
		fieldsList = append(fieldsList, key)
		fieldValuesList = append(fieldValuesList, fmt.Sprint(value))
	}
	fieldsList = append(fieldsList, "message")
	fieldValuesList = append(fieldValuesList, el.defaults.Message)
	fieldsHeader := fmt.Sprintf(
		fieldsTemplate,
		strings.Join(fieldsList, " | "),
	)

	logLine := strings.Join(fieldValuesList, " | ")

	dateHeader := fmt.Sprintf(
		dateTemplate,
		el.defaults.Timestamp.String(),
	)

	outputString := fmt.Sprintf(
		formatTemplate,
		versionHeader,
		dateHeader,
		fieldsHeader,
		logLine,
	)

	return []byte(outputString)
}

func getELFLog(
	elfLogger *ELFLogger,
	level LogLevel,
	message string,
) *ELFLog {

	newELFLog := ELFLog{
		contents: make(map[string]interface{}),
		defaults: defaultValues{
			Message: message,
			Level:   level,
			Timestamp: Timestamp{
				Time: time.Now(),
			},
		},
		elfLogger: elfLogger,
	}

	return &newELFLog
}
