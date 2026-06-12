package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/user/chatbot-whatsapp/controllers"
)

type Handler struct {
	Webhook *controllers.WebhookController
	Message *controllers.MessageController
	Health  *controllers.HealthController
}

func Setup(router *gin.Engine, h *Handler, log *zap.Logger) {
	router.Use(corsMiddleware())

	router.GET("/health", h.Health.Check)

	webhook := router.Group("/webhook")
	{
		webhook.GET("", h.Webhook.HandleVerification)
		webhook.POST("", h.Webhook.HandleMessage)
	}

	api := router.Group("/api/v1")
	{
		api.POST("/messages/text", h.Message.SendText)
		api.POST("/messages/image", h.Message.SendImage)
		api.POST("/messages/template", h.Message.SendTemplate)
		api.POST("/messages/queue", h.Message.QueueMessage)
		api.POST("/broadcast", h.Message.Broadcast)
	}

	log.Info("routes configured",
		zap.String("webhook", "/webhook"),
		zap.String("api", "/api/v1"),
	)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
