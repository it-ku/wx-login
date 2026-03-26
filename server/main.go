package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v5"
	echomw "github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"wx-login/server/config"
	"wx-login/server/handler"
	"wx-login/server/middleware"
	"wx-login/server/model"
	"wx-login/server/service"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(
		&model.User{},
		&model.VerifyCodeLog{},
		&model.RateLimitLock{},
	); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	// 初始化服务
	codeSvc := service.NewCodeService(db)
	wechatSvc := service.NewWechatService(db, codeSvc, cfg.ServiceQrcodeMediaID)
	authSvc := service.NewAuthService(db, codeSvc, cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// 初始化 handler
	wechatHandler := handler.NewWechatHandler(cfg, wechatSvc)
	authHandler := handler.NewAuthHandler(authSvc)

	// 初始化 Echo
	e := echo.New()

	e.Use(echomw.Recover())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// 静态文件服务（二维码图片等）
	e.Static("/static", "static")

	// 路由
	api := e.Group("/api")

	// 微信回调（无需认证）
	api.GET("/wechat/callback", wechatHandler.Callback)
	api.POST("/wechat/callback", wechatHandler.Callback)

	// 认证接口（无需认证）
	api.POST("/auth/login", authHandler.Login)

	// 需要 JWT 的接口
	auth := api.Group("", middleware.JWTAuth(authSvc))
	auth.GET("/auth/user", authHandler.User)
	auth.POST("/auth/logout", authHandler.Logout)

	// 启动
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
