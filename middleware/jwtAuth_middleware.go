package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rmine/ktMicro/util/jwtAuth"
	ktMicro "github.com/rmine/ktMicro/util/responseUtil"
)

//token格式: cookie/header ==>  Authorization: Bearer <token>
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//校验client里的token
		_, claims, err := jwtAuth.AutoVerifyToken(c)
		if err != nil {
			ktMicro.ResponseJsonWithCode(c, nil, 1001)
			c.Abort()
			return
		}
		//获取用户信息
		authInfo, err := jwtAuth.GetValidAuthInfo(claims)
		if err != nil {
			ktMicro.ResponseJsonWithCode(c, nil, 1001)
			c.Abort()
			return
		}
		//插入用户信息进上下文, 无值不阻塞
		if authInfo != nil {
			if authInfo.UserId > 0 {
				//c.Set(pconst.AuthUserId, authInfo.UserId)
			}
			//if authInfo.UserInfo != nil {
			//	c.Set(pconst.AuthUserInfo,authInfo.UserInfo)
			//}

			//user, err := mysql.NewUserDao().FindUser(authInfo.UserId)
			//if err == nil {
			//	c.Set(pconst.AuthUserId,user)
			//}
		}

		c.Next()
	}
}
