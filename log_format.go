package logging

import (
)

// LogFormat is a representation of what format a log should output.
type LogFormat int

// Definition of all the available LogFormats known by this framework.
const (
    UnsetFormat LogFormat = iota
    JSON
    Standard
    StandardExtended
)
