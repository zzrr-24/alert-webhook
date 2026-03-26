package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type FeishuService interface {
	SendCard(card map[string]interface{}) error
}

type feishuService struct {
	webhookURL string
	timeout    time.Duration
	logger     *zap.Logger
}

func NewFeishuService(webhookURL string, timeout time.Duration, logger *zap.Logger) FeishuService {
	return &feishuService{
		webhookURL: webhookURL,
		timeout:    timeout,
		logger:     logger,
	}
}

func (s *feishuService) SendCard(card map[string]interface{}) error {
	cardBytes, err := json.Marshal(card)
	if err != nil {
		s.logger.Error("marshal card error", zap.Error(err))
		return fmt.Errorf("marshal card error: %w", err)
	}

	s.logger.Info("feishu request body", zap.String("body", string(cardBytes)))

	client := &http.Client{Timeout: s.timeout}
	resp, err := client.Post(s.webhookURL, "application/json", bytes.NewBuffer(cardBytes))
	if err != nil {
		s.logger.Error("send error", zap.Error(err))
		return fmt.Errorf("send error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	s.logger.Info("feishu response", zap.Int("status", resp.StatusCode), zap.String("body", string(respBody)))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("feishu api error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	return nil
}
