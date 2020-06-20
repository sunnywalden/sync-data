package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetLogger, init a logger
func GetLogger(logLevel logrus.Level) (logger *logrus.Logger) {
	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logLevel

	logger.WithFields(logrus.Fields{
		"app_id": "user",
	})

	return logger
}

// Logger, logger for router
func Logger(logLevel logrus.Level) gin.HandlerFunc {
		// use logrus
		logger := logrus.New()

		logger.Formatter = &logrus.JSONFormatter{}
		//设置日志格式

		logger.SetFormatter(
			&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logger.Level = logLevel
	return func(c *gin.Context) {

		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"latency_time": latencyTime,
			"client_ip": clientIP,
			"req_method": reqMethod,
			"req_uri": reqUri,
		}).Info()
	}
}
