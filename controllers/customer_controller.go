package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/user/chatbot-whatsapp/models"
)

type CustomerController struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCustomerController(db *gorm.DB, log *zap.Logger) *CustomerController {
	return &CustomerController{db: db, log: log}
}

func (ctrl *CustomerController) List(c *gin.Context) {
	if ctrl.db == nil {
		c.JSON(http.StatusOK, gin.H{"customers": []models.Customer{}})
		return
	}

	var customers []models.Customer
	if err := ctrl.db.Order("created_at DESC").Limit(100).Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list customers"})
		return
	}

	if customers == nil {
		customers = []models.Customer{}
	}

	c.JSON(http.StatusOK, gin.H{"customers": customers})
}
