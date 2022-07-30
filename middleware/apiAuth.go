package middleware

import (
	"github.com/gin-gonic/gin"
	ktMicro "github.com/rmine/ktMicro/util/responseUtil"
)

func ApiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok := checkUserLogin(c)
		if ok {
			c.Next()
		} else {
			ktMicro.ResponseJsonWithCode(c, nil, 1001)
			c.Abort()
		}
	}
}

func checkUserLogin(c *gin.Context) (ok bool) {
	username, _ := c.Cookie("username")
	if len(username) == 0 {
		return false
	}

	token, _ := c.Cookie("token")
	if len(token) == 0 {
		return false
	}

	//ok, err := mysql.NewUserTokenDao().IsValidToken(username, token)
	//if err != nil {
	//	return false
	//}
	c.Set("username", username)
	return
}
