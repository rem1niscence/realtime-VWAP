package logger

import (
	"github.com/sirupsen/logrus"
)

var logLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

// Log is a custom instance of the logrus logger
var Log = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	Log.SetFormatter(&logrus.JSONFormatter{})
}

// SetLogLevel sets the level of logs, defaults to info
func SetLogLevel(level string) {
	lvl, ok := logLevels[level]
	if !ok {
		Log.SetLevel(logrus.InfoLevel)
	}
	Log.SetLevel(lvl)
}
