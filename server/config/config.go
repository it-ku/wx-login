package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Wechat  WechatConfig  `yaml:"wechat"`
	JWT     JWTConfig     `yaml:"jwt"`
	Database DatabaseConfig `yaml:"database"`

	WechatQrcode          string `yaml:"wechat_qrcode"`
	ServiceQrcode         string `yaml:"service_qrcode"`
	ServiceQrcodeMediaID  string `yaml:"service_qrcode_media_id"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type WechatConfig struct {
	AppID          string `yaml:"app_id"`
	AppSecret      string `yaml:"app_secret"`
	Token          string `yaml:"token"`
	EncodingAESKey string `yaml:"encoding_aes_key"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.JWT.ExpireHours == 0 {
		cfg.JWT.ExpireHours = 24
	}
	return &cfg, nil
}
