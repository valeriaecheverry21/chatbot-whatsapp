package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/user/chatbot-whatsapp/config"
	"github.com/user/chatbot-whatsapp/models"
	"github.com/user/chatbot-whatsapp/queue"
	"github.com/user/chatbot-whatsapp/services"
)

type MessageController struct {
	cfg             *config.Config
	whatsAppService *services.WhatsAppService
	customerService *services.CustomerService
	messageService  *services.MessageService
	queueClient     *queue.Client
	log             *zap.Logger
}

func NewMessageController(
	cfg *config.Config,
	whatsAppService *services.WhatsAppService,
	customerService *services.CustomerService,
	messageService *services.MessageService,
	queueClient *queue.Client,
	log *zap.Logger,
) *MessageController {
	return &MessageController{
		cfg:             cfg,
		whatsAppService: whatsAppService,
		customerService: customerService,
		messageService:  messageService,
		queueClient:     queueClient,
		log:             log,
	}
}

type SendTextRequest struct {
	Phone   string `json:"phone" binding:"required"`
	Message string `json:"message" binding:"required"`
	Preview bool   `json:"preview"`
}

type SendImageRequest struct {
	Phone   string `json:"phone" binding:"required"`
	URL     string `json:"url" binding:"required"`
	Caption string `json:"caption"`
}

type SendTemplateRequest struct {
	Phone    string   `json:"phone" binding:"required"`
	Name     string   `json:"name" binding:"required"`
	Language string   `json:"language"`
	Params   []string `json:"params"`
}

type BroadcastRequest struct {
	TemplateName string   `json:"template_name" binding:"required"`
	Language     string   `json:"language"`
	Params       []string `json:"params"`
	Phones       []string `json:"phones,omitempty"`
}

func (ctrl *MessageController) SendText(c *gin.Context) {
	var req SendTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := ctrl.whatsAppService.SendText(c.Request.Context(), services.SendTextInput{
		Phone:   req.Phone,
		Message: req.Message,
		Preview: req.Preview,
	})
	if err != nil {
		ctrl.log.Error("send text failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}

	if err := ctrl.messageService.Save(msg); err != nil {
		ctrl.log.Error("failed to save message", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "sent",
		"whatsapp_id": msg.WhatsAppID,
	})
}

func (ctrl *MessageController) SendImage(c *gin.Context) {
	var req SendImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := ctrl.whatsAppService.SendImage(c.Request.Context(), services.SendImageInput{
		Phone:   req.Phone,
		URL:     req.URL,
		Caption: req.Caption,
	})
	if err != nil {
		ctrl.log.Error("send image failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send image"})
		return
	}

	if err := ctrl.messageService.Save(msg); err != nil {
		ctrl.log.Error("failed to save message", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "sent",
		"whatsapp_id": msg.WhatsAppID,
	})
}

func (ctrl *MessageController) SendTemplate(c *gin.Context) {
	var req SendTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lang := req.Language
	if lang == "" {
		lang = "es"
	}

	msg, err := ctrl.whatsAppService.SendTemplate(c.Request.Context(), services.SendTemplateInput{
		Phone:    req.Phone,
		Name:     req.Name,
		Language: lang,
		Params:   req.Params,
	})
	if err != nil {
		ctrl.log.Error("send template failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send template"})
		return
	}

	if err := ctrl.messageService.Save(msg); err != nil {
		ctrl.log.Error("failed to save message", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "sent",
		"whatsapp_id": msg.WhatsAppID,
	})
}

func (ctrl *MessageController) QueueMessage(c *gin.Context) {
	var req SendTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &queue.QueueMessage{
		Phone:       req.Phone,
		MessageType: models.MessageTypeText,
		Body:        req.Message,
		CreatedAt:   models.Message{}.CreatedAt,
	}

	if err := ctrl.queueClient.Publish(c.Request.Context(), msg); err != nil {
		ctrl.log.Error("queue publish failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue message"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "queued",
		"phone":  req.Phone,
	})
}

func (ctrl *MessageController) Broadcast(c *gin.Context) {
	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lang := req.Language
	if lang == "" {
		lang = "es"
	}

	var customers []models.Customer
	var err error

	if len(req.Phones) > 0 {
		for _, phone := range req.Phones {
			customer, findErr := ctrl.customerService.FindByPhone(phone)
			if findErr == nil && customer.OptIn {
				customers = append(customers, *customer)
			}
		}
	} else {
		customers, err = ctrl.customerService.ListOptedIn(1000, 0)
		if err != nil {
			ctrl.log.Error("failed to list customers", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list customers"})
			return
		}
	}

	queued := 0
	for _, customer := range customers {
		if !customer.OptIn {
			continue
		}

		queueMsg := &queue.QueueMessage{
			CustomerID:  customer.ID,
			Phone:       customer.Phone,
			MessageType: models.MessageTypeTemplate,
			Template: &queue.QueueTemplate{
				Name:     req.TemplateName,
				Language: lang,
				Params:   req.Params,
			},
		}

		if err := ctrl.queueClient.Publish(c.Request.Context(), queueMsg); err != nil {
			ctrl.log.Error("failed to queue broadcast message",
				zap.String("phone", customer.Phone),
				zap.Error(err),
			)
			continue
		}
		queued++
	}

	ctrl.log.Info("broadcast queued",
		zap.Int("total_customers", len(customers)),
		zap.Int("queued", queued),
	)

	c.JSON(http.StatusAccepted, gin.H{
		"status":            "broadcast_queued",
		"total_customers":   len(customers),
		"messages_queued":   queued,
	})
}
