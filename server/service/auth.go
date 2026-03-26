package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"wx-login/server/model"
)

var (
	ErrCodeInvalid = errors.New("验证码无效或已过期")
)

// AuthService 认证服务
type AuthService struct {
	db          *gorm.DB
	codeSvc     *CodeService
	jwtSecret   []byte
	expireHours int
}

func NewAuthService(db *gorm.DB, codeSvc *CodeService, jwtSecret string, expireHours int) *AuthService {
	return &AuthService{
		db:          db,
		codeSvc:     codeSvc,
		jwtSecret:   []byte(jwtSecret),
		expireHours: expireHours,
	}
}

// Login 验证码登录，返回 token、expiresAt 和 User
func (s *AuthService) Login(code string) (token string, expiresAt time.Time, user *model.User, err error) {
	openID, ok := s.codeSvc.Verify(code)
	if !ok {
		err = ErrCodeInvalid
		return
	}

	// 创建或查找用户
	var u model.User
	result := s.db.Where("open_id = ?", openID).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		u = model.User{OpenID: openID, Subscribed: true}
		if err = s.db.Create(&u).Error; err != nil {
			return
		}
	} else if result.Error != nil {
		err = result.Error
		return
	}

	// 签发 JWT
	expiresAt = time.Now().Add(time.Duration(s.expireHours) * time.Hour)
	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%d", u.ID),
		"oid": openID,
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString(s.jwtSecret)
	if err != nil {
		return
	}
	user = &u
	return
}

// GetUserByOpenID 通过 OpenID 获取用户
func (s *AuthService) GetUserByOpenID(openID string) (*model.User, error) {
	var u model.User
	if err := s.db.Where("open_id = ?", openID).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// ParseToken 解析 JWT，返回 openID
func (s *AuthService) ParseToken(tokenStr string) (openID string, err error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !t.Valid {
		return "", errors.New("token invalid")
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	openID, _ = claims["oid"].(string)
	return openID, nil
}
