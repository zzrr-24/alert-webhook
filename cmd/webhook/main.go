package main

import (
	"fmt"
	"log"

	"alert-webhook/internal/config"
	"alert-webhook/internal/handler"
	"alert-webhook/internal/service"
	loggerpkg "alert-webhook/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}

	if err := loggerpkg.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		log.Fatalf("init logger error: %v", err)
	}
	defer loggerpkg.Logger.Sync()

	loggerpkg.Info("Feishu adapter starting", zap.String("webhook", cfg.Feishu.WebhookURL))

	cardService := service.NewCardService()
	feishuService := service.NewFeishuService(cfg.Feishu.WebhookURL, cfg.Feishu.Timeout, loggerpkg.Logger)
	handler := handler.NewHandler(loggerpkg.Logger, cardService, feishuService)

	gin.SetMode(cfg.Server.Mode)
	router := gin.New()
	router.Use(gin.Recovery())
	handler.RegisterRoutes(router)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	loggerpkg.Info("Feishu adapter running", zap.String("addr", addr))

	if err := router.Run(addr); err != nil {
		loggerpkg.Error("server error", zap.Error(err))
		log.Fatalf("server error: %v", err)
	}
}
