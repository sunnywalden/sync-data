package logging

import (

	"github.com/sirupsen/logrus"

	"github.com/sunnywalden/sync-data/config"
)


// GetLogger, init a logger
func GetLogger(configures *config.LogConf) (logger *logrus.Logger) {

	logLevel := configures.Level

	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logLevel

	return logger
}
