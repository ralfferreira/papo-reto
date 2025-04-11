package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/models"
	"gorm.io/gorm"
)

// MessageRepository handles database operations for messages
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

// Create creates a new message
func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

// GetByID gets a message by ID
func (r *MessageRepository) GetByID(id uuid.UUID) (*models.Message, error) {
	var message models.Message
	if err := r.db.First(&message, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}
	return &message, nil
}

// GetByGroupID gets all messages for a group
func (r *MessageRepository) GetByGroupID(groupID uuid.UUID) ([]models.Message, error) {
	var messages []models.Message
	if err := r.db.Where("group_id = ?", groupID).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// GetByGroupIDPaginated gets paginated messages for a group
func (r *MessageRepository) GetByGroupIDPaginated(groupID uuid.UUID, page, pageSize int) ([]models.Message, error) {
	var messages []models.Message
	offset := (page - 1) * pageSize
	if err := r.db.Where("group_id = ?", groupID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// Update updates a message
func (r *MessageRepository) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

// Delete deletes a message
func (r *MessageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Message{}, "id = ?", id).Error
}

// MarkAsRead marks a message as read
func (r *MessageRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&models.Message{}).Where("id = ?", id).
		Update("is_read", true).Error
}

// ToggleFavorite toggles the favorite status of a message
func (r *MessageRepository) ToggleFavorite(id uuid.UUID) error {
	var message models.Message
	if err := r.db.First(&message, "id = ?", id).Error; err != nil {
		return err
	}

	return r.db.Model(&message).Update("is_favorite", !message.IsFavorite).Error
}

// RevealIdentity reveals the identity of a message sender
func (r *MessageRepository) RevealIdentity(id uuid.UUID, senderID string) error {
	return r.db.Model(&models.Message{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_revealed": true,
			"sender_id":   senderID,
		}).Error
}

// CountByGroupID counts the number of messages in a group
func (r *MessageRepository) CountByGroupID(groupID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Message{}).
		Where("group_id = ?", groupID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByUserIDInPeriod counts the number of messages for a user in a specific period
func (r *MessageRepository) CountByUserIDInPeriod(userID uuid.UUID, startTime, endTime time.Time) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Message{}).
		Joins("JOIN message_groups ON messages.group_id = message_groups.id").
		Where("message_groups.user_id = ? AND messages.created_at BETWEEN ? AND ?", userID, startTime, endTime).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// AnonymizeOldIPs anonymizes IP addresses for messages older than the specified duration
func (r *MessageRepository) AnonymizeOldIPs(olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)
	return r.db.Model(&models.Message{}).
		Where("created_at < ? AND sender_ip != 'anonymized'", cutoffTime).
		Update("sender_ip", "anonymized").Error
}
