package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"alert-webhook/internal/model"
	"alert-webhook/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger        *zap.Logger
	cardService   service.CardService
	feishuService service.FeishuService
}

func NewHandler(logger *zap.Logger, cardService service.CardService, feishuService service.FeishuService) *Handler {
	return &Handler{
		logger:        logger,
		cardService:   cardService,
		feishuService: feishuService,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/alert", h.HandleAlert)
}

func (h *Handler) HandleAlert(c *gin.Context) {
	h.logger.Info("==== Received Alertmanager Webhook ====")

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("read body error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body error"})
		return
	}

	h.logger.Info("raw body", zap.String("body", string(bodyBytes)))

	var payload model.AlertmanagerWebhook
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&payload); err != nil {
		h.logger.Error("decode error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "decode error"})
		return
	}

	h.logger.Info("parsed alerts", zap.Int("count", len(payload.Alerts)), zap.String("status", payload.Status))

	card := h.cardService.BuildCard(payload)

	if err := h.feishuService.SendCard(card); err != nil {
		h.logger.Error("send card error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "send card error"})
		return
	}

	h.logger.Info("==== Alert processed successfully ====")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
