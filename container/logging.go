package container

import "github.com/sirupsen/logrus"

func ConfigureLogger(logLevel logrus.Level) *logrus.Logger {
	log := logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)
	log.SetReportCaller(true)
	log.SetLevel(logLevel)
	return log
}

var logger = ConfigureLogger(logrus.InfoLevel)
