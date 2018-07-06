package slogging

// LogInstance is an instance of a log entry.
type LogInstance interface {
	And(string, interface{}) LogInstance
	With(string, interface{}) LogInstance
	Send()
	SetTimestamp(Timestamp) LogInstance
}
