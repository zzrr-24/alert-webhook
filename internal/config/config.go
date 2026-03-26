package config

import "time"

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Feishu FeishuConfig `mapstructure:"feishu"`
	Log    LogConfig    `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type FeishuConfig struct {
	WebhookURL string        `mapstructure:"webhook_url"`
	Timeout    time.Duration `mapstructure:"timeout"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}
