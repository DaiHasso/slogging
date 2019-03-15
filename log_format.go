package logging

import (
    "strings"
)

// LogFormat is a representation of what format a log should output.
type LogFormat int

// Definition of all the available LogFormats known by this framework.
const (
    UnsetFormat LogFormat = iota
    UnknownFormat
    JSON
    Standard
    StandardExtended
)

func FormatFromString(format string) LogFormat {
    switch strings.ToLower(format) {
    case "json":
        return JSON
    case "standard":
        return Standard
    case "standardextended":
        return StandardExtended
    default:
        return UnknownFormat
    }
}
