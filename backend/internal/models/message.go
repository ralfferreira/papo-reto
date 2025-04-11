package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message represents an anonymous message sent to a group
type Message struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	GroupID    uuid.UUID `gorm:"type:uuid;index"`
	Content    string    `gorm:"type:text"`
	SenderIP   string    `gorm:"size:50"`
	SenderID   *string   `gorm:"size:255"` // Optional, for revealed identity
	IsRead     bool      `gorm:"default:false"`
	IsFavorite bool      `gorm:"default:false"`
	IsRevealed bool      `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`

	// Define this as a belongs-to relationship with the correct references
	Group MessageGroup `gorm:"foreignKey:GroupID;references:ID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// MarkAsRead marks the message as read
func (m *Message) MarkAsRead() {
	m.IsRead = true
}

// ToggleFavorite toggles the favorite status of the message
func (m *Message) ToggleFavorite() {
	m.IsFavorite = !m.IsFavorite
}

// RevealIdentity marks the message as having a revealed identity
func (m *Message) RevealIdentity(senderID string) {
	m.IsRevealed = true
	m.SenderID = &senderID
}

// AnonymizeIP anonymizes the sender's IP address for privacy
func (m *Message) AnonymizeIP() {
	// Replace the last octet with zeros for IPv4 or truncate IPv6
	m.SenderIP = "anonymized"
}
