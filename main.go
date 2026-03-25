package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type AlertmanagerWebhook struct {
	Status string `json:"status"`
	Alerts []struct {
		Status      string            `json:"status"`
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
		StartsAt    time.Time         `json:"startsAt"`
	} `json:"alerts"`
}

var feishuWebhook = os.Getenv("FEISHU_WEBHOOK")

func main() {
	http.HandleFunc("/alert", handleAlert)
	log.Println("FEISHU_WEBHOOK:", feishuWebhook)
	log.Println("Feishu adapter running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAlert(w http.ResponseWriter, r *http.Request) {
	log.Println("==== Received Alertmanager Webhook ====")

	// 读取原始 body（用于日志）
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read body error:", err)
		http.Error(w, err.Error(), 400)
		return
	}
	log.Println("raw body:", string(bodyBytes))

	// 重新放回 body 给 decoder 用
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload AlertmanagerWebhook
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("decode error:", err)
		http.Error(w, err.Error(), 400)
		return
	}

	log.Printf("parsed alerts: %d, status: %s\n", len(payload.Alerts), payload.Status)

	card := buildCard(payload)

	cardBytes, _ := json.Marshal(card)
	log.Println("feishu request body:", string(cardBytes))

	resp, err := http.Post(feishuWebhook, "application/json", bytes.NewBuffer(cardBytes))
	if err != nil {
		log.Println("send error:", err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("feishu response status: %d, body: %s\n", resp.StatusCode, string(respBody))

	log.Println("==== Alert processed successfully ====")

	w.Write([]byte("ok"))
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

func statusEmoji(status string) string {
	if status == "firing" {
		return "🔥"
	}
	return "✅"
}

func buildCard(p AlertmanagerWebhook) map[string]interface{} {
	headerColor := "blue"

	if len(p.Alerts) > 0 {
		headerColor = severityColor(p.Alerts[0].Labels["severity"])
	}

	var elements []map[string]interface{}

	for i, a := range p.Alerts {
		severity := a.Labels["severity"]
		statusIcon := statusEmoji(a.Status)
		sevIcon := severityEmoji(severity)

		title := fmt.Sprintf("%s %s %s", statusIcon, sevIcon, a.Labels["alertname"])
		summary := a.Annotations["summary"]

		// 基础信息块
		info := fmt.Sprintf(
			"**状态**: %s\n**级别**: %s\n**实例**: `%s`\n**开始时间**: %s",
			strings.ToUpper(a.Status),
			severity,
			a.Labels["instance"],
			a.StartsAt.Format("2006-01-02 15:04:05"),
		)

		// 标签整理
		var labelList []string
		for k, v := range a.Labels {
			if k == "alertname" || k == "severity" {
				continue
			}
			labelList = append(labelList, fmt.Sprintf("`%s=%s`", k, v))
		}

		labels := strings.Join(labelList, " ")

		// 卡片内容
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

		// 标签折叠（可选）
		if labels != "" {
			elements = append(elements, map[string]interface{}{
				"tag": "div",
				"text": map[string]string{
					"tag":     "lark_md",
					"content": fmt.Sprintf("**标签**: %s", labels),
				},
			})
		}

		// 分割线
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
					"content": fmt.Sprintf("🚨 Prometheus 告警（%d 条）", len(p.Alerts)),
				},
			},
			"elements": elements,
		},
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
