package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestSendCard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0}`))
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	service := NewFeishuService(server.URL, 10*time.Second, logger)

	card := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"header": map[string]string{
				"title": "test",
			},
		},
	}

	err := service.SendCard(card)
	assert.NoError(t, err)
}

func TestSendCardError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code":500,"msg":"internal error"}`))
	}))
	defer server.Close()

	logger, _ := zap.NewDevelopment()
	service := NewFeishuService(server.URL, 10*time.Second, logger)

	card := map[string]interface{}{
		"msg_type": "interactive",
	}

	err := service.SendCard(card)
	assert.Error(t, err)
}
