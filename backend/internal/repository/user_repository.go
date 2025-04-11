package repository

import (
	"encoding/json"
	"errors"
	"time"

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

// ListUsers lists all users with pagination
func (r *UserRepository) ListUsers(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * pageSize

	// Count total records
	r.db.Model(&models.User{}).Count(&total)

	// Fetch users with pagination
	result := r.db.Offset(offset).Limit(pageSize).Find(&users)
	return users, total, result.Error
}

// ListPremiumUsers lists all premium users
func (r *UserRepository) ListPremiumUsers() ([]models.User, error) {
	var users []models.User
	result := r.db.Where("plan = ?", "premium").Find(&users)
	return users, result.Error
}

// GetRecentUsers gets users registered recently
func (r *UserRepository) GetRecentUsers(days int) ([]models.User, error) {
	var users []models.User
	cutoffDate := time.Now().AddDate(0, 0, -days)

	result := r.db.Where("created_at >= ?", cutoffDate).Find(&users)
	return users, result.Error
}

// SearchUsers searches users by name or email
func (r *UserRepository) SearchUsers(query string) ([]models.User, error) {
	var users []models.User
	searchQuery := "%" + query + "%"

	result := r.db.Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery).Find(&users)
	return users, result.Error
}

// GetUserWithStats gets a user with updated usage statistics
func (r *UserRepository) GetUserWithStats(id uuid.UUID) (*models.User, error) {
	var user models.User

	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	// Update user statistics
	var activeGroups int64
	r.db.Model(&models.MessageGroup{}).Where("user_id = ? AND is_archived = false", id).Count(&activeGroups)
	user.ActiveGroups = int(activeGroups)

	// Update message count
	var messageCount int64
	r.db.Model(&models.Message{}).
		Joins("JOIN message_groups ON messages.group_id = message_groups.id").
		Where("message_groups.user_id = ?", id).
		Count(&messageCount)
	user.MessageCount = int(messageCount)

	// Save updated statistics
	r.db.Save(&user)

	return &user, nil
}

// GetUsersExceedingLimit finds users who have exceeded their plan's message limit
func (r *UserRepository) GetUsersExceedingLimit() ([]models.User, error) {
	var users []models.User

	// This query finds free plan users who have exceeded the 50 message limit
	result := r.db.Where("plan = ? AND message_count > ?", "free", 50).Find(&users)
	return users, result.Error
}

// GetUsersByVerificationStatus gets users by verification status
func (r *UserRepository) GetUsersByVerificationStatus(verified bool) ([]models.User, error) {
	var users []models.User
	result := r.db.Where("is_verified = ?", verified).Find(&users)
	return users, result.Error
}

// UpdateNotificationSettings updates a user's notification settings
func (r *UserRepository) UpdateNotificationSettings(userID uuid.UUID, settings json.RawMessage) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("notify_settings", settings).Error
}

// GetUsersForBulkEmail gets users eligible to receive bulk emails
func (r *UserRepository) GetUsersForBulkEmail() ([]models.User, error) {
	var users []models.User

	// Assuming notification settings include an option for marketing emails
	result := r.db.Where("is_verified = ? AND notify_settings->'email_marketing' = 'true'", true).Find(&users)
	return users, result.Error
}

// CountUsersByPlan counts the number of users by plan
func (r *UserRepository) CountUsersByPlan() (map[string]int64, error) {
	var results []struct {
		Plan  string
		Count int64
	}

	err := r.db.Model(&models.User{}).
		Select("plan, count(*) as count").
		Group("plan").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int64)
	for _, result := range results {
		countMap[result.Plan] = result.Count
	}

	return countMap, nil
}

// GetUserRegistrationStats gets user registration statistics by period
func (r *UserRepository) GetUserRegistrationStats(interval string) ([]struct {
	Period string
	Count  int64
}, error) {
	var results []struct {
		Period string
		Count  int64
	}

	var timeFormat string
	switch interval {
	case "day":
		timeFormat = "YYYY-MM-DD"
	case "week":
		timeFormat = "YYYY-WW"
	case "month":
		timeFormat = "YYYY-MM"
	default:
		timeFormat = "YYYY-MM-DD"
	}

	err := r.db.Model(&models.User{}).
		Select("TO_CHAR(created_at, ?) as period, count(*) as count", timeFormat).
		Group("period").
		Order("period").
		Find(&results).Error

	return results, err
}
