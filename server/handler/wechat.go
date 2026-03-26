package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"wx-login/server/config"
	"wx-login/server/service"
)

// WechatHandler 微信消息处理
type WechatHandler struct {
	cfg       *config.Config
	wechatSvc *service.WechatService
}

func NewWechatHandler(cfg *config.Config, wechatSvc *service.WechatService) *WechatHandler {
	return &WechatHandler{
		cfg:       cfg,
		wechatSvc: wechatSvc,
	}
}

// Callback 处理微信推送（GET: URL验证 / POST: 消息事件）
func (h *WechatHandler) Callback(c *echo.Context) error {
	wc := wechat.NewWechat()
	oa := wc.GetOfficialAccount(&offConfig.Config{
		AppID:          h.cfg.Wechat.AppID,
		AppSecret:      h.cfg.Wechat.AppSecret,
		Token:          h.cfg.Wechat.Token,
		EncodingAESKey: h.cfg.Wechat.EncodingAESKey,
		Cache:          cache.NewMemory(),
	})

	srv := oa.GetServer(c.Request(), c.Response())

	srv.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		openID := string(msg.FromUserName)

		switch msg.MsgType {
		case message.MsgTypeEvent:
			switch msg.Event {
			case message.EventSubscribe:
				return h.wechatSvc.HandleSubscribe(openID)
			case message.EventUnsubscribe:
				h.wechatSvc.HandleUnsubscribe(openID)
				return nil
			}

		case message.MsgTypeText:
			switch string(msg.Content) {
			case "验证码":
				return h.wechatSvc.HandleVerifyCode(openID)
			case "联系客服":
				return h.wechatSvc.HandleServiceQrcode()
			default:
				return h.wechatSvc.HandleDefault()
			}

		default:
			return h.wechatSvc.HandleDefault()
		}
		return nil
	})

	if err := srv.Serve(); err != nil {
		return c.String(http.StatusForbidden, err.Error())
	}
	if err := srv.Send(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return nil
}
