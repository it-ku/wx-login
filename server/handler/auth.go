package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"wx-login/server/middleware"
	"wx-login/server/service"
)

// AuthHandler 登录认证
type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type loginRequest struct {
	Code string `json:"code"`
}

// Login POST /api/auth/login
func (h *AuthHandler) Login(c *echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil || len(req.Code) != 6 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    40000,
			"message": "请输入6位数字验证码",
		})
	}

	token, expiresAt, user, err := h.authSvc.Login(req.Code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    40001,
			"message": "验证码无效或已过期",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "登录成功",
		"data": map[string]interface{}{
			"token":      token,
			"expires_at": expiresAt,
			"user":       user,
		},
	})
}

// User GET /api/auth/user
func (h *AuthHandler) User(c *echo.Context) error {
	openID := c.Get(middleware.ContextKeyOpenID).(string)
	user, err := h.authSvc.GetUserByOpenID(openID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    40100,
			"message": "未登录或 Token 已过期",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

// Logout POST /api/auth/logout
func (h *AuthHandler) Logout(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "退出成功",
	})
}
