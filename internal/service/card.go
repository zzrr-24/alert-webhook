package service

import (
	"fmt"
	"strings"
	"time"

	"alert-webhook/internal/model"
)

type CardService interface {
	BuildCard(p model.AlertmanagerWebhook) map[string]interface{}
}

type cardService struct{}

func NewCardService() CardService {
	return &cardService{}
}

func (s *cardService) BuildCard(p model.AlertmanagerWebhook) map[string]interface{} {
	var titlePrefix string
	var sevIcon string
	isResolved := p.Status == "resolved"

	headerColor := "blue"

	if isResolved {
		headerColor = "green"
		titlePrefix = "✅ Prometheus 告警恢复"
		sevIcon = "🟢"
	} else if len(p.Alerts) > 0 {
		headerColor = severityColor(p.Alerts[0].Labels["severity"])
		titlePrefix = "🚨 Prometheus 告警"
	}

	var elements []map[string]interface{}

	for i, a := range p.Alerts {
		severity := a.Labels["severity"]
		statusIcon := statusEmoji(a.Status)
		if sevIcon != "🟢" {
			sevIcon = severityEmoji(severity)
		}

		title := fmt.Sprintf("%s %s %s", statusIcon, sevIcon, a.Labels["alertname"])
		summary := a.Annotations["summary"]

		info := fmt.Sprintf(
			"**状态**: %s\n**级别**: %s\n**实例**: `%s`\n**开始时间**: %s",
			strings.ToUpper(a.Status),
			severity,
			a.Labels["instance"],
			formatTime(a.StartsAt),
		)

		var labelList []string
		for k, v := range a.Labels {
			if k == "alertname" || k == "severity" {
				continue
			}
			labelList = append(labelList, fmt.Sprintf("`%s=%s`", k, v))
		}

		labels := strings.Join(labelList, " ")

		elements = append(elements,
			map[string]interface{}{
				"tag": "div",
				"text": map[string]string{
					"tag":     "lark_md",
					"content": fmt.Sprintf("%s\n%s", title, summary),
				},
			},
			map[string]interface{}{
				"tag": "div",
				"text": map[string]string{
					"tag":     "lark_md",
					"content": info,
				},
			},
		)

		if labels != "" {
			elements = append(elements, map[string]interface{}{
				"tag": "div",
				"text": map[string]string{
					"tag":     "lark_md",
					"content": fmt.Sprintf("**标签**: %s", labels),
				},
			})
		}

		if i < len(p.Alerts)-1 {
			elements = append(elements, map[string]interface{}{
				"tag": "hr",
			})
		}
	}

	return map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"config": map[string]interface{}{
				"wide_screen_mode": true,
			},
			"header": map[string]interface{}{
				"template": headerColor,
				"title": map[string]string{
					"tag":     "plain_text",
					"content": fmt.Sprintf("%s（%d 条）", titlePrefix, len(p.Alerts)),
				},
			},
			"elements": elements,
		},
	}
}

func formatTime(t time.Time) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return t.In(loc).Format("2006-01-02 15:04:05")
}

func severityEmoji(sev string) string {
	switch sev {
	case "critical":
		return "🔴"
	case "warning":
		return "🟠"
	default:
		return "🔵"
	}
}

func severityColor(sev string) string {
	switch sev {
	case "critical":
		return "red"
	case "warning":
		return "orange"
	default:
		return "blue"
	}
}

func statusEmoji(status string) string {
	if status == "firing" {
		return "🔥"
	}
	return "✅"
}
