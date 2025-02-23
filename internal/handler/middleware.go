package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RootRequired 检查是否为管理员用户的中间件
func RootRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "未登录",
			})
			c.Abort()
			return
		}

		if role.(string) != "admin" {
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "权限不足，需要管理员权限",
				})
			} else {
				c.Redirect(http.StatusFound, "/")
			}
			c.Abort()
			return
		}

		c.Next()
	}
}
