package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// IncrementMessageCount increments the message count for a user
func (r *UserRepository) IncrementMessageCount(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("message_count", gorm.Expr("message_count + ?", 1)).Error
}

// IncrementActiveGroups increments the active groups count for a user
func (r *UserRepository) IncrementActiveGroups(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("active_groups", gorm.Expr("active_groups + ?", 1)).Error
}

// DecrementActiveGroups decrements the active groups count for a user
func (r *UserRepository) DecrementActiveGroups(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("active_groups", gorm.Expr("active_groups - ?", 1)).Error
}

// UpdatePlan updates the user's plan
func (r *UserRepository) UpdatePlan(userID uuid.UUID, plan string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("plan", plan).Error
}
