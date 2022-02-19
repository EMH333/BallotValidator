package util

import "time"

// Raw is the raw row
// Timestamp is the timestamp of the row (parsed)
// ONID is the ONID of the row (taken from the raw)
type Vote struct {
	Raw       []string
	Timestamp time.Time
	ONID      string
	ID        string
}

type Summary struct {
	Processed int
	Valid     int
	Invalid   int
	Log       []string
	StepInfo  string
}
