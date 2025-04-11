package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a registered user in the system
type User struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key"`
	Email          string          `gorm:"size:255;uniqueIndex"`
	Password       string          `gorm:"size:255"`
	Name           string          `gorm:"size:100"`
	AvatarURL      string          `gorm:"size:255"`
	IsVerified     bool            `gorm:"default:false"`
	Plan           string          `gorm:"size:50;default:'free'"`
	MessageCount   int             `gorm:"default:0"`
	ActiveGroups   int             `gorm:"default:0"`
	NotifySettings json.RawMessage `gorm:"type:jsonb"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// IsPremium checks if the user has a premium plan
func (u *User) IsPremium() bool {
	return u.Plan == "premium"
}

// GetGroupLimit returns the maximum number of active groups allowed for the user's plan
func (u *User) GetGroupLimit() int {
	if u.IsPremium() {
		return -1 // Unlimited
	}
	return 3 // Free plan limit
}

// GetMessageLimit returns the maximum number of messages allowed per month for the user's plan
func (u *User) GetMessageLimit() int {
	if u.IsPremium() {
		return -1 // Unlimited
	}
	return 50 // Free plan limit
}
