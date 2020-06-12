package logging

import (

	"github.com/sirupsen/logrus"
)


func GetLogger() (logger *logrus.Logger) {

	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.DebugLevel

	return logger
}
