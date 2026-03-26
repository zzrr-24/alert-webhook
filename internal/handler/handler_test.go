package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"alert-webhook/internal/model"
	"alert-webhook/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockFeishuService struct {
	sendError error
}

func (m *mockFeishuService) SendCard(card map[string]interface{}) error {
	return m.sendError
}

func TestHandleAlert(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	cardService := service.NewCardService()
	feishuService := &mockFeishuService{sendError: nil}

	handler := NewHandler(logger, cardService, feishuService)
	router := gin.New()
	handler.RegisterRoutes(router)

	now := time.Now()
	payload := model.AlertmanagerWebhook{
		Status: "firing",
		Alerts: []model.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "TestAlert",
					"severity":  "warning",
					"instance":  "localhost:8080",
				},
				Annotations: map[string]string{
					"summary": "Test summary",
				},
				StartsAt: now,
			},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/alert", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleAlertSendError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	cardService := service.NewCardService()
	feishuService := &mockFeishuService{sendError: assert.AnError}

	handler := NewHandler(logger, cardService, feishuService)
	router := gin.New()
	handler.RegisterRoutes(router)

	now := time.Now()
	payload := model.AlertmanagerWebhook{
		Status: "firing",
		Alerts: []model.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "TestAlert",
					"severity":  "warning",
					"instance":  "localhost:8080",
				},
				Annotations: map[string]string{
					"summary": "Test summary",
				},
				StartsAt: now,
			},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/alert", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
