package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/user/chatbot-whatsapp/queue"
)

type HealthController struct {
	db          *gorm.DB
	queueClient *queue.Client
	log         *zap.Logger
}

func NewHealthController(db *gorm.DB, queueClient *queue.Client, log *zap.Logger) *HealthController {
	return &HealthController{db: db, queueClient: queueClient, log: log}
}

func (ctrl *HealthController) Check(c *gin.Context) {
	dbOK := true
	if ctrl.db != nil {
		if sqlDB, err := ctrl.db.DB(); err != nil || sqlDB.Ping() != nil {
			dbOK = false
		}
	} else {
		dbOK = false
	}

	redisOK := true
	if ctrl.queueClient != nil {
		if err := ctrl.queueClient.Ping(c.Request.Context()); err != nil {
			redisOK = false
		}
	} else {
		redisOK = false
	}

	status := http.StatusOK
	if !dbOK || !redisOK {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"status":  "ok",
		"service": "whatsapp-chatbot",
		"checks": gin.H{
			"database": dbOK,
			"redis":    redisOK,
		},
	})
}
