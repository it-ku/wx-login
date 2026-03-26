package model

import "time"

// VerifyCodeLog 验证码发送记录表
type VerifyCodeLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OpenID    string    `gorm:"index;size:128;comment:微信OpenID" json:"open_id"`
	Code      string    `gorm:"size:6;comment:验证码" json:"code"`
	CreatedAt time.Time `json:"created_at"`
}
