package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"wx-login/server/service"
)

const ContextKeyOpenID = "openID"

// JWTAuth JWT 认证中间件
func JWTAuth(authSvc *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    40100,
					"message": "未登录或 Token 已过期",
				})
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			openID, err := authSvc.ParseToken(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    40100,
					"message": "未登录或 Token 已过期",
				})
			}
			c.Set(ContextKeyOpenID, openID)
			return next(c)
		}
	}
}
