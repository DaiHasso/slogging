package logging

import (
	"fmt"
	"time"
)

// timestamp is a time.Time that marshals to unix timestamp.
type timestamp struct {
	time.Time
}

// MarshalJSON will marshal the timestamp into a unix timestamp.
func (t timestamp) MarshalJSON() ([]byte, error) {
	unixTimestamp := fmt.Sprintf("%v", t.Unix())
	return []byte(unixTimestamp), nil
}

func (t timestamp) String() string {
    timeBytes := t.Format("2006-01-02T15:04:05")
	return string(timeBytes)
}
