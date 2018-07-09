package slogging

import (
	"encoding/json"
	"strings"
	"time"
)

type defaultValues struct {
	Message   string    `json:"message"`
	Level     LogLevel  `json:"log_level"`
	Timestamp timestamp `json:"timestamp"`
}

// JSONLog is an instance of a JSON log entry.
type JSONLog struct {
	jsonLogger *JSONLogger
	contents   map[string]interface{}
	defaults   defaultValues
	pretty     bool
}

// MarshalJSON will format the JSON log entry to be flat.
func (jl JSONLog) MarshalJSON() ([]byte, error) {
	var outputJSON []byte

	outputJSON, err := json.Marshal(jl.defaults)
	if err != nil {
		return nil, err
	}
	if len(jl.contents) == 0 {
		return outputJSON, nil
	}

	contentsJSON, err := json.Marshal(jl.contents)
	if err != nil {
		return nil, err
	}
	outputJSON[len(outputJSON)-1] = ','
	outputJSON = append(outputJSON, contentsJSON[1:]...)

	return outputJSON, nil
}

// SetTimestamp overrides timestamp, mostly used for testing.
func (jl *JSONLog) SetTimestamp(newTime time.Time) LogInstance {
	jl.defaults.Timestamp = timestamp{newTime}
	return jl
}

// And performs the same functionality as @With.
func (jl *JSONLog) And(key string, val interface{}) LogInstance {
	return jl.With(key, val)
}

// With adds a key and value to a JSONLog.
func (jl *JSONLog) With(key string, val interface{}) LogInstance {
	jl.contents[strings.ToLower(key)] = val
	return jl
}

// Send outputs the log entry to the logger and ensures proper
// formatting.
func (jl *JSONLog) Send() {
	var marshaledOutput []byte
	if jl.pretty {
		marshaledOutput, _ = json.MarshalIndent(jl, "", "  ")
	} else {
		marshaledOutput, _ = json.Marshal(jl)
	}
	jl.jsonLogger.Log(
		jl.defaults.Level,
		marshaledOutput,
	)
}

// GetJSONLog initializes and returns a new JSON log.
func getJSONLog(
	jsonLogger *JSONLogger,
	level LogLevel,
	message string,
	pretty bool,
) *JSONLog {

	newJSONLog := JSONLog{
		contents: make(map[string]interface{}),
		defaults: defaultValues{
			Message: message,
			Level:   level,
			Timestamp: timestamp{
				Time: time.Now(),
			},
		},
		jsonLogger: jsonLogger,
		pretty:     pretty,
	}
	return &newJSONLog
}
