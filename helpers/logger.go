package helpers

import "github.com/sirupsen/logrus"

// SetLogLevel sets the log level for logrus by string; possible values are debug, error, fatal, panic, info, trace
func SetLogLevel(loglevelString string) {
	loglevel := logrus.ErrorLevel
	switch loglevelString {
	case "":
		loglevel = logrus.ErrorLevel
	case "error":
		loglevel = logrus.ErrorLevel
	case "debug":
		loglevel = logrus.DebugLevel
	case "fatal":
		loglevel = logrus.FatalLevel
	case "panic":
		loglevel = logrus.PanicLevel
	case "info":
		loglevel = logrus.InfoLevel
	case "trace":
		loglevel = logrus.TraceLevel
	default:
		loglevel = logrus.ErrorLevel
	}
	logrus.SetLevel(loglevel)
	logrus.Debugf("LogLevel: %s", loglevel.String())

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		FullTimestamp:             true,
		QuoteEmptyFields:          true,
		EnvironmentOverrideColors: true,
	})
}
