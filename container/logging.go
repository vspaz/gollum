package container

import "github.com/sirupsen/logrus"

func ConfigureLogger(logLevel logrus.Level) *logrus.Logger {
	logger := logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"
	formatter.FullTimestamp = true
	logger.SetFormatter(formatter)
	logger.SetReportCaller(true)
	logger.SetLevel(logLevel)
	return logger
}
