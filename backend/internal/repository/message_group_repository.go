package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/models"
	"gorm.io/gorm"
)

// MessageGroupRepository handles database operations for message groups
type MessageGroupRepository struct {
	db *gorm.DB
}

// NewMessageGroupRepository creates a new message group repository
func NewMessageGroupRepository(db *gorm.DB) *MessageGroupRepository {
	return &MessageGroupRepository{
		db: db,
	}
}

// Create creates a new message group
func (r *MessageGroupRepository) Create(group *models.MessageGroup) error {
	return r.db.Create(group).Error
}

// GetByID gets a message group by ID
func (r *MessageGroupRepository) GetByID(id uuid.UUID) (*models.MessageGroup, error) {
	var group models.MessageGroup
	if err := r.db.First(&group, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("message group not found")
		}
		return nil, err
	}
	return &group, nil
}

// GetBySlug gets a message group by slug
func (r *MessageGroupRepository) GetBySlug(slug string) (*models.MessageGroup, error) {
	var group models.MessageGroup
	if err := r.db.First(&group, "slug = ?", slug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("message group not found")
		}
		return nil, err
	}
	return &group, nil
}

// GetByUserID gets all message groups for a user
func (r *MessageGroupRepository) GetByUserID(userID uuid.UUID) ([]models.MessageGroup, error) {
	var groups []models.MessageGroup
	if err := r.db.Where("user_id = ?", userID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetActiveByUserID gets all active message groups for a user
func (r *MessageGroupRepository) GetActiveByUserID(userID uuid.UUID) ([]models.MessageGroup, error) {
	var groups []models.MessageGroup
	if err := r.db.Where("user_id = ? AND is_archived = ?", userID, false).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// Update updates a message group
func (r *MessageGroupRepository) Update(group *models.MessageGroup) error {
	return r.db.Save(group).Error
}

// Archive archives a message group
func (r *MessageGroupRepository) Archive(id uuid.UUID) error {
	return r.db.Model(&models.MessageGroup{}).Where("id = ?", id).
		Update("is_archived", true).Error
}

// Unarchive unarchives a message group
func (r *MessageGroupRepository) Unarchive(id uuid.UUID) error {
	return r.db.Model(&models.MessageGroup{}).Where("id = ?", id).
		Update("is_archived", false).Error
}

// Delete deletes a message group
func (r *MessageGroupRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.MessageGroup{}, "id = ?", id).Error
}

// CountActiveByUserID counts the number of active groups for a user
func (r *MessageGroupRepository) CountActiveByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.MessageGroup{}).
		Where("user_id = ? AND is_archived = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// IsSlugAvailable checks if a slug is available
func (r *MessageGroupRepository) IsSlugAvailable(slug string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.MessageGroup{}).
		Where("slug = ?", slug).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
