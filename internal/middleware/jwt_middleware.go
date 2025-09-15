package middleware

import (
	"goblog/pkg/response"
	"goblog/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.从请求头获取Authorization
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, 401, "请先登录")
			c.Abort()
			return
		}
		//2.检查Authorization格式，必须是bearer+空格+token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, 401, "令牌格式错误")
			c.Abort()
			return
		}

		//3.调用jwt工具包解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 401, "令牌已过期或无效："+err.Error())
			c.Abort()
			return
		}
		//4.将用户信息存储上下文
		c.Set("userID", claims.ID)
		c.Next()
	}
}
