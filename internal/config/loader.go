package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Load(configPath string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("feishu.timeout", "10s")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}

	v.SetEnvPrefix("ALERT_WEBHOOK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}

	return &cfg, nil
}
