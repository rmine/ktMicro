package middleware

import (
	"github.com/gin-gonic/gin"
	ktMicro "github.com/rmine/ktMicro/config"
	"github.com/sirupsen/logrus"
	"time"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//请求参数
		postParams := c.Request.PostForm.Encode()
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := c.ClientIP()

		go ktMicro.GinLogger().WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
			"post_params":  postParams,
		}).Info()
	}
}
