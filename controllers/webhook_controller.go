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

type WebhookController struct {
	cfg             *config.WhatsAppConfig
	whatsAppService *services.WhatsAppService
	customerService *services.CustomerService
	messageService  *services.MessageService
	queueClient     *queue.Client
	log             *zap.Logger
}

func NewWebhookController(
	cfg *config.WhatsAppConfig,
	whatsAppService *services.WhatsAppService,
	customerService *services.CustomerService,
	messageService *services.MessageService,
	queueClient *queue.Client,
	log *zap.Logger,
) *WebhookController {
	return &WebhookController{
		cfg:             cfg,
		whatsAppService: whatsAppService,
		customerService: customerService,
		messageService:  messageService,
		queueClient:     queueClient,
		log:             log,
	}
}

type WebhookEntry struct {
	Messages []InboundMessage `json:"messages,omitempty"`
	Statuses []StatusUpdate   `json:"statuses,omitempty"`
}

type InboundMessage struct {
	ID       string   `json:"id"`
	From     string   `json:"from"`
	Type     string   `json:"type"`
	Text     *TextMsg `json:"text,omitempty"`
	Image    *Media   `json:"image,omitempty"`
	Document *Media   `json:"document,omitempty"`
	Audio    *Media   `json:"audio,omitempty"`
	Video    *Media   `json:"video,omitempty"`
}

type TextMsg struct {
	Body string `json:"body"`
}

type Media struct {
	ID       string `json:"id,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	Caption  string `json:"caption,omitempty"`
}

type StatusUpdate struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type WebhookCallback struct {
	Object string         `json:"object"`
	Entry  []WebhookEntry `json:"entry"`
}

func (ctrl *WebhookController) HandleVerification(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	ctrl.log.Debug("webhook verification request",
		zap.String("mode", mode),
		zap.String("token", token),
	)

	if mode == "subscribe" && token == ctrl.cfg.VerifyToken {
		ctrl.log.Info("webhook verified successfully")
		c.String(http.StatusOK, challenge)
		return
	}

	ctrl.log.Warn("webhook verification failed",
		zap.String("expected_token", ctrl.cfg.VerifyToken),
	)
	c.AbortWithStatus(http.StatusForbidden)
}

func (ctrl *WebhookController) HandleMessage(c *gin.Context) {
	var callback WebhookCallback
	if err := c.ShouldBindJSON(&callback); err != nil {
		ctrl.log.Error("invalid webhook payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	ctrl.log.Debug("webhook received",
		zap.String("object", callback.Object),
		zap.Int("entries", len(callback.Entry)),
	)

	for _, entry := range callback.Entry {
		for _, msg := range entry.Messages {
			ctrl.processInboundMessage(c, msg)
		}
		for _, status := range entry.Statuses {
			ctrl.processStatusUpdate(c, status)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (ctrl *WebhookController) processInboundMessage(c *gin.Context, msg InboundMessage) {
	ctrl.log.Info("inbound message",
		zap.String("from", msg.From),
		zap.String("type", msg.Type),
		zap.String("id", msg.ID),
	)

	customer, err := ctrl.customerService.FindOrCreate(msg.From)
	if err != nil {
		ctrl.log.Error("failed to find/create customer", zap.Error(err))
		return
	}

	messageType := mapMessageType(msg.Type)

	record := &models.Message{
		WhatsAppID:  msg.ID,
		CustomerID:  customer.ID,
		Direction:   models.DirectionInbound,
		MessageType: messageType,
		Status:      models.StatusDelivered,
		Body:        extractBody(msg),
	}

	if err := ctrl.messageService.Save(record); err != nil {
		ctrl.log.Error("failed to save inbound message", zap.Error(err))
	}

	ctrl.whatsAppService.MarkAsRead(c.Request.Context(), msg.ID)
}

func (ctrl *WebhookController) processStatusUpdate(c *gin.Context, status StatusUpdate) {
	ctrl.log.Debug("status update",
		zap.String("id", status.ID),
		zap.String("status", status.Status),
	)

	msgStatus := mapStatus(status.Status)
	if msgStatus == "" {
		return
	}

	if err := ctrl.messageService.UpdateStatus(status.ID, msgStatus); err != nil {
		ctrl.log.Warn("failed to update message status",
			zap.String("id", status.ID),
			zap.Error(err),
		)
	}
}

func extractBody(msg InboundMessage) string {
	switch {
	case msg.Text != nil:
		return msg.Text.Body
	case msg.Image != nil:
		return msg.Image.Caption
	case msg.Document != nil:
		return msg.Document.Caption
	default:
		return ""
	}
}

func mapMessageType(t string) models.MessageType {
	switch t {
	case "text":
		return models.MessageTypeText
	case "image":
		return models.MessageTypeImage
	case "document":
		return models.MessageTypeDocument
	case "audio":
		return models.MessageTypeAudio
	case "video":
		return models.MessageTypeVideo
	case "sticker":
		return models.MessageTypeSticker
	case "location":
		return models.MessageTypeLocation
	case "contacts":
		return models.MessageTypeContacts
	case "button":
		return models.MessageTypeButton
	case "interactive":
		return models.MessageTypeInterac
	default:
		return models.MessageTypeText
	}
}

func mapStatus(s string) models.MessageStatus {
	switch s {
	case "sent":
		return models.StatusSent
	case "delivered":
		return models.StatusDelivered
	case "read":
		return models.StatusRead
	case "failed":
		return models.StatusFailed
	default:
		return ""
	}
}
