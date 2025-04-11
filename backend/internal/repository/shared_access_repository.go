package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/models"
	"gorm.io/gorm"
)

// SharedAccessRepository handles database operations for shared access
type SharedAccessRepository struct {
	db *gorm.DB
}

// NewSharedAccessRepository creates a new shared access repository
func NewSharedAccessRepository(db *gorm.DB) *SharedAccessRepository {
	return &SharedAccessRepository{
		db: db,
	}
}

// Create creates a new shared access
func (r *SharedAccessRepository) Create(access *models.SharedAccess) error {
	return r.db.Create(access).Error
}

// GetByID gets a shared access by ID
func (r *SharedAccessRepository) GetByID(id uuid.UUID) (*models.SharedAccess, error) {
	var access models.SharedAccess
	if err := r.db.First(&access, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shared access not found")
		}
		return nil, err
	}
	return &access, nil
}

// GetByToken gets a shared access by token
func (r *SharedAccessRepository) GetByToken(token string) (*models.SharedAccess, error) {
	var access models.SharedAccess
	if err := r.db.First(&access, "token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shared access not found")
		}
		return nil, err
	}
	return &access, nil
}

// GetByGroupID gets all shared access for a group
func (r *SharedAccessRepository) GetByGroupID(groupID uuid.UUID) ([]models.SharedAccess, error) {
	var accesses []models.SharedAccess
	if err := r.db.Where("group_id = ?", groupID).Find(&accesses).Error; err != nil {
		return nil, err
	}
	return accesses, nil
}

// GetActiveByGroupID gets all active shared access for a group
func (r *SharedAccessRepository) GetActiveByGroupID(groupID uuid.UUID) ([]models.SharedAccess, error) {
	var accesses []models.SharedAccess
	now := time.Now()
	if err := r.db.Where("group_id = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)",
		groupID, true, now).Find(&accesses).Error; err != nil {
		return nil, err
	}
	return accesses, nil
}

// Update updates a shared access
func (r *SharedAccessRepository) Update(access *models.SharedAccess) error {
	return r.db.Save(access).Error
}

// Revoke revokes a shared access
func (r *SharedAccessRepository) Revoke(id uuid.UUID) error {
	return r.db.Model(&models.SharedAccess{}).Where("id = ?", id).
		Update("is_active", false).Error
}

// RevokeByGroupID revokes all shared access for a group
func (r *SharedAccessRepository) RevokeByGroupID(groupID uuid.UUID) error {
	return r.db.Model(&models.SharedAccess{}).Where("group_id = ?", groupID).
		Update("is_active", false).Error
}

// Delete deletes a shared access
func (r *SharedAccessRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.SharedAccess{}, "id = ?", id).Error
}

// CountByGroupID counts the number of active shared access for a group
func (r *SharedAccessRepository) CountByGroupID(groupID uuid.UUID) (int64, error) {
	var count int64
	now := time.Now()
	if err := r.db.Model(&models.SharedAccess{}).
		Where("group_id = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)",
			groupID, true, now).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CleanupExpired deletes all expired shared access
func (r *SharedAccessRepository) CleanupExpired() error {
	now := time.Now()
	return r.db.Where("expires_at < ?", now).Delete(&models.SharedAccess{}).Error
}
