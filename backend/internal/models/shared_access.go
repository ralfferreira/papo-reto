package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SharedAccess represents shared access to a message group
type SharedAccess struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key"`
	GroupID   uuid.UUID  `gorm:"type:uuid;index"`
	InvitedBy uuid.UUID  `gorm:"type:uuid"`
	Email     string     `gorm:"size:255"`
	Token     string     `gorm:"size:100;uniqueIndex"`
	IsActive  bool       `gorm:"default:true"`
	ExpiresAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	
	Group     MessageGroup `gorm:"foreignKey:GroupID"`
	Inviter   User         `gorm:"foreignKey:InvitedBy"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (sa *SharedAccess) BeforeCreate(tx *gorm.DB) error {
	if sa.ID == uuid.Nil {
		sa.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the shared access has expired
func (sa *SharedAccess) IsExpired() bool {
	if sa.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*sa.ExpiresAt)
}

// IsValid checks if the shared access is valid (active and not expired)
func (sa *SharedAccess) IsValid() bool {
	return sa.IsActive && !sa.IsExpired()
}

// Revoke deactivates the shared access
func (sa *SharedAccess) Revoke() {
	sa.IsActive = false
}
