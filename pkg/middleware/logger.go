package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ko3-gin/pkg/logger"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		allTime := endTime.Sub(startTime)
		reqMethod := ctx.Request.Method
		reqUrl := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		logger.Default.WithFields(logrus.Fields{
			"all_time":  allTime,
			"client_ip": clientIP,
		}).Debug(fmt.Sprintf("request: [%s] %s   code: %d", reqMethod, reqUrl, statusCode))
	}
}
