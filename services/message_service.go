package services

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/user/chatbot-whatsapp/models"
)

type MessageService struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewMessageService(db *gorm.DB, log *zap.Logger) *MessageService {
	return &MessageService{db: db, log: log}
}

func (s *MessageService) Save(msg *models.Message) error {
	return s.db.Create(msg).Error
}

func (s *MessageService) UpdateStatus(whatsAppID string, status models.MessageStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	switch status {
	case models.StatusSent:
		updates["sent_at"] = gorm.Expr("NOW()")
	case models.StatusDelivered:
		updates["delivered_at"] = gorm.Expr("NOW()")
	case models.StatusRead:
		updates["read_at"] = gorm.Expr("NOW()")
	}

	return s.db.Model(&models.Message{}).
		Where("whatsapp_id = ?", whatsAppID).
		Updates(updates).Error
}

func (s *MessageService) UpdateError(whatsAppID, errorCode string) error {
	return s.db.Model(&models.Message{}).
		Where("whatsapp_id = ?", whatsAppID).
		Updates(map[string]interface{}{
			"status":     models.StatusFailed,
			"error_code": errorCode,
		}).Error
}

func (s *MessageService) FindByWhatsAppID(whatsAppID string) (*models.Message, error) {
	var msg models.Message
	err := s.db.Where("whatsapp_id = ?", whatsAppID).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (s *MessageService) ListByCustomer(customerID uint, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := s.db.Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (s *MessageService) ListByStatus(status models.MessageStatus, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := s.db.Where("status = ?", status).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}
