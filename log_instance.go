package logging

import "time"

// LogInstance is an instance of a log entry.
type LogInstance interface {
	And(string, interface{}) LogInstance
	With(string, interface{}) LogInstance
	Send()
	SetTimestamp(time.Time) LogInstance
}
