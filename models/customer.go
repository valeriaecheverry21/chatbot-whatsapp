package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Phone     string         `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Name      string         `gorm:"size:255" json:"name"`
	OptIn     bool           `gorm:"default:false;not null" json:"opt_in"`
	OptInAt   *time.Time     `json:"opt_in_at,omitempty"`
	Language  string         `gorm:"size:10;default:es" json:"language"`
	Tags      string         `gorm:"type:text" json:"tags,omitempty"`
	Notes     string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
