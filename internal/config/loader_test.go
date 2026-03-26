package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefaults(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Feishu: FeishuConfig{
			WebhookURL: "https://open.feishu.cn/open-apis/bot/v2/hook/test",
			Timeout:    10000000000,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.Mode)
	assert.Equal(t, "https://open.feishu.cn/open-apis/bot/v2/hook/test", cfg.Feishu.WebhookURL)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
}
