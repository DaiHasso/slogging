package slogging

// PseudoWriter is a wrapper for JSONLogger for things that
// need a writer to output.
type PseudoWriter struct {
	logger   Logger
	logLevel LogLevel
}

func (pw PseudoWriter) Write(p []byte) (n int, err error) {
	switch pw.logLevel {
	case ERROR:
		pw.logger.Error(string(p)).Send()
	case WARN:
		pw.logger.Warn(string(p)).Send()
	case INFO:
		pw.logger.Info(string(p)).Send()
	case DEBUG:
		pw.logger.Debug(string(p)).Send()
	}
	return len(p), nil
}
