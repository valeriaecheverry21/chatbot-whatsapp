package models

import (
	"time"
)

type MessageDirection string

const (
	DirectionInbound  MessageDirection = "inbound"
	DirectionOutbound MessageDirection = "outbound"
)

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeTemplate MessageType = "template"
	MessageTypeDocument MessageType = "document"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeVideo    MessageType = "video"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeLocation MessageType = "location"
	MessageTypeContacts MessageType = "contacts"
	MessageTypeButton   MessageType = "button"
	MessageTypeInterac  MessageType = "interactive"
)

type MessageStatus string

const (
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
	StatusFailed    MessageStatus = "failed"
	StatusPending   MessageStatus = "pending"
)

type Message struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	WhatsAppID      string          `gorm:"uniqueIndex;size:100" json:"whatsapp_id"`
	CustomerID      uint            `gorm:"index;not null" json:"customer_id"`
	Customer        Customer        `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Direction       MessageDirection `gorm:"size:10;not null;index" json:"direction"`
	MessageType     MessageType     `gorm:"size:20;not null" json:"message_type"`
	Status          MessageStatus   `gorm:"size:20;default:pending;index" json:"status"`
	Body            string          `gorm:"type:text" json:"body"`
	MediaURL        string          `gorm:"size:500" json:"media_url,omitempty"`
	TemplateName    string          `gorm:"size:100" json:"template_name,omitempty"`
	TemplateLang    string          `gorm:"size:10" json:"template_language,omitempty"`
	TemplateParams  string          `gorm:"type:text" json:"template_params,omitempty"`
	ConversationID  string          `gorm:"size:100" json:"conversation_id,omitempty"`
	ErrorCode       string          `gorm:"size:50" json:"error_code,omitempty"`
	Metadata        string          `gorm:"type:jsonb" json:"metadata,omitempty"`
	SentAt          *time.Time      `json:"sent_at,omitempty"`
	DeliveredAt     *time.Time      `json:"delivered_at,omitempty"`
	ReadAt          *time.Time      `json:"read_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
