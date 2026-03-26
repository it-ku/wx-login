package model

import "time"

// RateLimitLock 频次锁定表
type RateLimitLock struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OpenID      string    `gorm:"uniqueIndex;size:128;comment:微信OpenID" json:"open_id"`
	LockedUntil time.Time `gorm:"comment:锁定截止时间" json:"locked_until"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
