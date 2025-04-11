package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MessageGroup represents a group of anonymous messages
type MessageGroup struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID       `gorm:"type:uuid"`
	Name        string          `gorm:"size:100"`
	Slug        string          `gorm:"size:100;uniqueIndex"`
	Description string          `gorm:"size:500"`
	IsPublic    bool            `gorm:"default:false"`
	IsArchived  bool            `gorm:"default:false"`
	Settings    json.RawMessage `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`

	User User `gorm:"foreignKey:UserID"`
	// Define this as a has-many relationship with the correct references
	Messages []Message `gorm:"foreignKey:GroupID;references:ID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (mg *MessageGroup) BeforeCreate(tx *gorm.DB) error {
	if mg.ID == uuid.Nil {
		mg.ID = uuid.New()
	}
	return nil
}

// IsActive returns whether the group is active (not archived)
func (mg *MessageGroup) IsActive() bool {
	return !mg.IsArchived
}

// GetIcebreakers returns the icebreaker questions configured for this group
func (mg *MessageGroup) GetIcebreakers() []string {
	if mg.Settings == nil {
		return []string{}
	}

	var settings struct {
		Icebreakers []string `json:"icebreakers"`
	}

	if err := json.Unmarshal(mg.Settings, &settings); err != nil {
		return []string{}
	}

	return settings.Icebreakers
}

// GetBannedWords returns the list of banned words for content moderation
func (mg *MessageGroup) GetBannedWords() []string {
	if mg.Settings == nil {
		return []string{}
	}

	var settings struct {
		BannedWords []string `json:"bannedWords"`
	}

	if err := json.Unmarshal(mg.Settings, &settings); err != nil {
		return []string{}
	}

	return settings.BannedWords
}
