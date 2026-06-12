package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/user/chatbot-whatsapp/config"
	"github.com/user/chatbot-whatsapp/controllers"
	"github.com/user/chatbot-whatsapp/database"
	"github.com/user/chatbot-whatsapp/queue"
	"github.com/user/chatbot-whatsapp/routes"
	"github.com/user/chatbot-whatsapp/services"
	"github.com/user/chatbot-whatsapp/workers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := config.NewLogger(cfg)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("starting whatsapp chatbot service",
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)

	db, err := database.Connect(&cfg.Database, logger)
	if err != nil {
		logger.Warn("database unavailable, running without persistence", zap.Error(err))
	} else {
		if err := database.Migrate(db, logger); err != nil {
			logger.Warn("database migration failed", zap.Error(err))
		}
	}

	queueClient := queue.NewClient(&cfg.Redis, logger)
	if err := queueClient.Ping(context.Background()); err != nil {
		logger.Warn("redis unavailable, running without queue", zap.Error(err))
		queueClient = nil
	}
	if queueClient != nil {
		defer queueClient.Close()
	}

	whatsAppService := services.NewWhatsAppService(&cfg.WhatsApp, logger)
	customerService := services.NewCustomerService(db, logger)
	messageService := services.NewMessageService(db, logger)

	webhookCtrl := controllers.NewWebhookController(
		&cfg.WhatsApp, whatsAppService, customerService, messageService, queueClient, logger,
	)
	messageCtrl := controllers.NewMessageController(
		cfg, whatsAppService, customerService, messageService, queueClient, logger,
	)
	healthCtrl := controllers.NewHealthController(db, queueClient, logger)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(logger))

	routes.Setup(router, &routes.Handler{
		Webhook: webhookCtrl,
		Message: messageCtrl,
		Health:  healthCtrl,
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if queueClient != nil {
		workerPool := workers.NewWorkerPool(5, queueClient, whatsAppService, messageService, logger)
		go workerPool.Start(ctx)
	} else {
		logger.Warn("worker pool not started: no queue available")
	}

	srv := &server{
		router: router,
		port:   cfg.Server.Port,
		log:    logger,
	}

	go srv.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("shutting down server", zap.String("signal", sig.String()))
	cancel()
	srv.Stop()
	logger.Info("server stopped gracefully")
}

type server struct {
	router *gin.Engine
	port   string
	log    *zap.Logger
}

func (s *server) Start() {
	s.log.Info("http server listening", zap.String("port", s.port))
	if err := s.router.Run(":" + s.port); err != nil {
		s.log.Fatal("server error", zap.Error(err))
	}
}

func (s *server) Stop() {
	s.log.Info("http server stopped")
}

func requestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
		)
	}
}
