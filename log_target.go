package slogging

// LogTarget is a representation of where a log should output to.
//go:generate stringer -type=LogTarget
type LogTarget int

// Definition of LogTargets for a logger.
const (
	_ LogTarget = iota
	Stdout
	LogFile
)
