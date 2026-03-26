package model

import "time"

// User 用户表
type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OpenID     string    `gorm:"uniqueIndex;size:128;comment:微信OpenID" json:"open_id"`
	Subscribed bool      `gorm:"default:true;comment:是否关注公众号" json:"subscribed"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
