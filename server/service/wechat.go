package service

import (
	"fmt"
	"time"

	"github.com/silenceper/wechat/v2/officialaccount/message"
	"gorm.io/gorm"
	"wx-login/server/model"
)

// WechatService 微信消息业务处理
type WechatService struct {
	db                   *gorm.DB
	codeSvc              *CodeService
	serviceQrcodeMediaID string
}

func NewWechatService(db *gorm.DB, codeSvc *CodeService, mediaID string) *WechatService {
	return &WechatService{
		db:                   db,
		codeSvc:              codeSvc,
		serviceQrcodeMediaID: mediaID,
	}
}

// HandleSubscribe 处理关注事件
func (s *WechatService) HandleSubscribe(openID string) *message.Reply {
	// 创建或更新用户
	var user model.User
	result := s.db.Where("open_id = ?", openID).First(&user)
	if result.Error != nil {
		s.db.Create(&model.User{OpenID: openID, Subscribed: true})
	} else {
		s.db.Model(&user).Update("subscribed", true)
	}
	return textReply(`欢迎关注！如需登录请发送"验证码"，如需帮助请发送"联系客服"。`)
}

// HandleUnsubscribe 处理取消关注事件
func (s *WechatService) HandleUnsubscribe(openID string) {
	s.db.Model(&model.User{}).Where("open_id = ?", openID).Update("subscribed", false)
}

// HandleVerifyCode 处理"验证码"关键词
func (s *WechatService) HandleVerifyCode(openID string) *message.Reply {
	locked, _ := s.codeSvc.CheckRateLimit(openID)
	if locked {
		return textReply("超过使用频次，请6小时后再试。")
	}

	code, expiresAt, err := s.codeSvc.Generate(openID)
	if err != nil {
		return textReply("系统繁忙，请稍后重试。")
	}

	now := time.Now()
	var timeStr string
	if expiresAt.Year() == now.Year() && expiresAt.Month() == now.Month() && expiresAt.Day() == now.Day() {
		timeStr = fmt.Sprintf("%02d:%02d", expiresAt.Hour(), expiresAt.Minute())
		return textReply(fmt.Sprintf("您的登录验证码为【%s】，请在 %s 前使用。请在网页中输入验证码完成登录。", code, timeStr))
	}
	timeStr = fmt.Sprintf("%d月%d号 %02d:%02d", int(expiresAt.Month()), expiresAt.Day(), expiresAt.Hour(), expiresAt.Minute())
	return textReply(fmt.Sprintf("您的登录验证码为【%s】，请在 %s 前使用。请在网页中输入验证码完成登录。", code, timeStr))
}

// HandleServiceQrcode 处理"联系客服"关键词
func (s *WechatService) HandleServiceQrcode() *message.Reply {
	if s.serviceQrcodeMediaID == "" {
		return textReply("请联系管理员获取客服二维码。")
	}
	return &message.Reply{
		MsgType: message.MsgTypeImage,
		MsgData: message.NewImage(s.serviceQrcodeMediaID),
	}
}

// HandleDefault 处理其他消息
func (s *WechatService) HandleDefault() *message.Reply {
	return textReply(`无法识别消息，如果需要登录请发送"验证码"，如需其他请联系客服`)
}

func textReply(content string) *message.Reply {
	return &message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText(content),
	}
}


