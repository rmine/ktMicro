package jwtAuth

import (
	"github.com/rmine/ktMicro/util/reformData"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

var apiAuthScecret string = "apiAuthScecret_ktMicro"
var authExpiresMin time.Duration = 15

type AuthInfo struct {
	UserId   uint32
	UserInfo interface{}
}

func NewAuthInfo(authScecret string, authExpiresMinSec time.Duration) *AuthInfo {
	m := &AuthInfo{}
	apiAuthScecret = authScecret
	authExpiresMin = authExpiresMinSec
	return m
}

func init() {

}

//创建令牌
func CreateToken(authInfo *AuthInfo) (data string, err error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = authInfo.UserId
	claims["authInfo"] = authInfo
	claims["exp"] = time.Now().Add(authExpiresMin * time.Minute).Unix()
	claims["user"] = authInfo.UserInfo
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	data, err = token.SignedString([]byte(apiAuthScecret))
	return
}

//取token,优先级cookie>header,  token标准格式 Authorization: Bearer <token>,解析需要空格分割
func ExtractToken(c *gin.Context) (data string, err error) {
	bearToken, err := c.Cookie("Authorization")
	if err != nil {
		data, err = getTokenStringFromHeader(c)
		return
	} else {
		strArr := strings.Split(bearToken, " ")
		if len(strArr) == 2 {
			data = strArr[1]
		} else {
			data, err = getTokenStringFromHeader(c)
		}
	}
	return
}

//验证token
func ParseToken(tokenString string) (token *jwt.Token, err error) {
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(apiAuthScecret), nil
	})
	if err != nil {
		return nil, err
	}
	return
}

//验证token是否有效
func VerifyToken(token *jwt.Token) (data jwt.MapClaims, err error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		data = claims
	} else {
		err = errors.New("invalid Authorization token!")
	}
	return
}

//获取用户信息
func GetValidAuthInfo(claims jwt.MapClaims) (data *AuthInfo, err error) {
	data = new(AuthInfo)
	err = reformData.ReformJsonToModel(claims["authInfo"], data)
	if err != nil {
		return
	}
	return
}

//一键校验
func AutoVerifyToken(c *gin.Context) (token *jwt.Token, claims jwt.MapClaims, err error) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		return
	}

	token, err = ParseToken(tokenString)
	if err != nil {
		return
	}

	claims, err = VerifyToken(token)
	if err != nil {
		return
	}

	return
}

//############# private methods
func getTokenStringFromHeader(c *gin.Context) (data string, err error) {
	bearToken := c.GetHeader("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		data = strArr[1]
	} else {
		err = errors.New("valid bearToken in Authorization!")
	}
	return data, err
}
