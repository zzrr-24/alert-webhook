package service

import (
	"testing"
	"time"

	"alert-webhook/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestBuildCard(t *testing.T) {
	service := NewCardService()

	now := time.Now()

	p := model.AlertmanagerWebhook{
		Status: "firing",
		Alerts: []model.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "HighCPU",
					"severity":  "critical",
					"instance":  "server1:8080",
				},
				Annotations: map[string]string{
					"summary": "CPU usage is above 90%",
				},
				StartsAt: now,
			},
		},
	}

	card := service.BuildCard(p)

	assert.Equal(t, "interactive", card["msg_type"])
	assert.NotNil(t, card["card"])

	cardData := card["card"].(map[string]interface{})
	assert.NotNil(t, cardData["header"])
	assert.NotNil(t, cardData["elements"])

	header := cardData["header"].(map[string]interface{})
	assert.Equal(t, "red", header["template"])
}

func TestSeverityEmoji(t *testing.T) {
	assert.Equal(t, "🔴", severityEmoji("critical"))
	assert.Equal(t, "🟠", severityEmoji("warning"))
	assert.Equal(t, "🔵", severityEmoji("info"))
}

func TestSeverityColor(t *testing.T) {
	assert.Equal(t, "red", severityColor("critical"))
	assert.Equal(t, "orange", severityColor("warning"))
	assert.Equal(t, "blue", severityColor("info"))
}

func TestStatusEmoji(t *testing.T) {
	assert.Equal(t, "🔥", statusEmoji("firing"))
	assert.Equal(t, "✅", statusEmoji("resolved"))
}
