package logging

// PseudoWriter is a wrapper for JSONLogger for things that
// need a writer to output.
type PseudoWriter struct {
	logger   *Logger
	logLevel LogLevel
}

// Write satisfies the io.Writer interface and writes the the logger it wraps.
func (self PseudoWriter) Write(p []byte) (n int, err error) {
	switch self.logLevel {
	case ERROR:
		self.logger.Error(string(p))
	case WARN:
		self.logger.Warn(string(p))
	case INFO:
		self.logger.Info(string(p))
	case DEBUG:
		self.logger.Debug(string(p))
	}
	return len(p), nil
}

// NewPsuedoWriter wraps a logger with the Write functionality wich writes out
// logs at a specified log level.
func NewPseudoWriter(logLevel LogLevel, logger *Logger) *PseudoWriter {
    return &PseudoWriter{
        logger: logger,
        logLevel: logLevel,
    }
}
