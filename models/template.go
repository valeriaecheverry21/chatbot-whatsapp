package models

import (
	"time"
)

type TemplateCategory string

const (
	TemplateMarketing   TemplateCategory = "MARKETING"
	TemplateUtility     TemplateCategory = "UTILITY"
	TemplateAuth        TemplateCategory = "AUTHENTICATION"
)

type TemplateStatus string

const (
	TemplateApproved  TemplateStatus = "APPROVED"
	TemplatePending   TemplateStatus = "PENDING"
	TemplateRejected  TemplateStatus = "REJECTED"
	TemplateDisabled  TemplateStatus = "DISABLED"
)

type Template struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Category    TemplateCategory `gorm:"size:20;not null" json:"category"`
	Status      TemplateStatus  `gorm:"size:20;default:PENDING" json:"status"`
	Language    string          `gorm:"size:10;default:es" json:"language"`
	BodyText    string          `gorm:"type:text;not null" json:"body_text"`
	HeaderType  string          `gorm:"size:20" json:"header_type,omitempty"`
	HeaderValue string          `gorm:"size:500" json:"header_value,omitempty"`
	FooterText  string          `gorm:"type:text" json:"footer_text,omitempty"`
	Variables   int             `gorm:"default:0" json:"variables"`
	MetaID      string          `gorm:"size:100" json:"meta_id,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
