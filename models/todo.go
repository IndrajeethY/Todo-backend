package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	ID              string     `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID          string     `gorm:"index:idx_user_title,unique;not null;type:char(36)" json:"user_id"`
	Title           string     `gorm:"not null;index:idx_user_title,unique" json:"title"`
	Description     string     `json:"description"`
	DueDate         *time.Time `json:"due_date"`
	Priority        string     `gorm:"type:text;default:'medium'" json:"priority"`
	Completed       bool       `gorm:"default:false" json:"completed"`
	NotifyEnabled   bool       `gorm:"default:true" json:"notify_enabled"`
	NotifyFrequency int        `gorm:"default:60" json:"notify_frequency"`
	OrderIndex      int        `gorm:"default:0" json:"order_index"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	TelegramEnabled bool       `gorm:"default:false" json:"telegram_enabled"`
	DiscordEnabled  bool       `gorm:"default:false" json:"discord_enabled"`
	NextNotifyAt    *time.Time `json:"next_notify_at"`
}

func (t *Todo) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	if t.UserID == "" {
		return errors.New("user_id required")
	}
	if t.Priority == "" {
		t.Priority = "medium"
	}
	return nil
}
