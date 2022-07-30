package ktMicro

import (
	"github.com/gin-gonic/gin"
	ktMicro "github.com/rmine/ktMicro/config"
	"net/http"
)

func ResponseJsonOK(c *gin.Context, data interface{}) {
	msg := ktMicro.GetMsgValue(0)
	handleResponse(c, data, 0, msg)
}

func ResponseJsonWithCode(c *gin.Context, data interface{}, code int) {
	msg := ktMicro.GetMsgValue(code)
	handleResponse(c, data, code, msg)
}

func handleResponse(c *gin.Context, data interface{}, code int, msg string) {
	var result gin.H
	if data == nil {
		result = gin.H{
			"code":    code,
			"message": msg,
		}
	} else {
		result = gin.H{
			"data":    data,
			"code":    code,
			"message": msg,
		}
	}

	if gin.Mode() == gin.TestMode {
		//todo record log
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
