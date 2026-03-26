package service

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
	"wx-login/server/model"
)

const (
	codeExpiry      = 10 * time.Minute
	rateLimitWindow = 5 * time.Minute
	rateLimitMax    = 5
	cooldownDur     = 6 * time.Hour
)

type codeEntry struct {
	OpenID    string
	ExpiresAt time.Time
}

// CodeService 验证码服务
type CodeService struct {
	db    *gorm.DB
	mu    sync.RWMutex
	cache map[string]*codeEntry // code -> entry
}

func NewCodeService(db *gorm.DB) *CodeService {
	return &CodeService{
		db:    db,
		cache: make(map[string]*codeEntry),
	}
}

// CheckRateLimit 检查频次限制。返回 (locked bool, lockedUntil time.Time)
func (s *CodeService) CheckRateLimit(openID string) (locked bool, lockedUntil time.Time) {
	// 先检查是否在冷却期
	var lock model.RateLimitLock
	if err := s.db.Where("open_id = ?", openID).First(&lock).Error; err == nil {
		if lock.LockedUntil.After(time.Now()) {
			return true, lock.LockedUntil
		}
		// 过期了，清理
		s.db.Delete(&lock)
	}

	// 检查最近 5 分钟内的请求次数
	var count int64
	since := time.Now().Add(-rateLimitWindow)
	s.db.Model(&model.VerifyCodeLog{}).
		Where("open_id = ? AND created_at >= ?", openID, since).
		Count(&count)

	if count >= rateLimitMax {
		// 写入冷却期
		until := time.Now().Add(cooldownDur)
		s.db.Save(&model.RateLimitLock{
			OpenID:      openID,
			LockedUntil: until,
		})
		return true, until
	}
	return false, time.Time{}
}

// Generate 生成验证码，写入缓存和数据库
func (s *CodeService) Generate(openID string) (code string, expiresAt time.Time, err error) {
	b := make([]byte, 3)
	if _, err = rand.Read(b); err != nil {
		return
	}
	n := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	code = fmt.Sprintf("%06d", n)

	expiresAt = time.Now().Add(codeExpiry)

	s.mu.Lock()
	s.cache[code] = &codeEntry{OpenID: openID, ExpiresAt: expiresAt}
	s.mu.Unlock()

	// 写发送记录
	s.db.Create(&model.VerifyCodeLog{OpenID: openID, Code: code})
	return
}

// Verify 验证验证码，返回绑定的 OpenID
func (s *CodeService) Verify(code string) (openID string, ok bool) {
	s.mu.RLock()
	entry, exists := s.cache[code]
	s.mu.RUnlock()
	if !exists || time.Now().After(entry.ExpiresAt) {
		return "", false
	}
	// 使用后删除
	s.mu.Lock()
	delete(s.cache, code)
	s.mu.Unlock()
	return entry.OpenID, true
}
