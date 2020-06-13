package logging

import (

	"github.com/sirupsen/logrus"

	"github.com/sunnywalden/sync-data/config"
)


func GetLogger() (logger *logrus.Logger) {

	configures := config.Conf
	logLevel := configures.Log.Level

	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logLevel

	return logger
}
