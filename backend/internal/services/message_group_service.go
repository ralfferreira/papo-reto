package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/models"
	"github.com/ralfferreira/papo-reto/internal/repository"
)

// MessageGroupService handles business logic for message groups
type MessageGroupService struct {
	groupRepo *repository.MessageGroupRepository
	userRepo  *repository.UserRepository
}

// NewMessageGroupService creates a new message group service
func NewMessageGroupService(groupRepo *repository.MessageGroupRepository, userRepo *repository.UserRepository) *MessageGroupService {
	return &MessageGroupService{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// CreateGroup creates a new message group
func (s *MessageGroupService) CreateGroup(userID uuid.UUID, name, description string, isPublic bool, settings map[string]interface{}) (*models.MessageGroup, error) {
	// Check if user can create a new group
	canCreate, err := s.canUserCreateGroup(userID)
	if err != nil {
		return nil, err
	}
	if !canCreate {
		return nil, errors.New("user has reached the maximum number of active groups")
	}

	// Generate a unique slug
	slug := s.generateSlug(name)

	// Convert settings to JSON
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	// Create group
	group := &models.MessageGroup{
		UserID:      userID,
		Name:        name,
		Slug:        slug,
		Description: description,
		IsPublic:    isPublic,
		IsArchived:  false,
		Settings:    settingsJSON,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save group
	if err := s.groupRepo.Create(group); err != nil {
		return nil, err
	}

	// Increment user's active groups count
	if err := s.userRepo.IncrementActiveGroups(userID); err != nil {
		return nil, err
	}

	return group, nil
}

// GetGroupByID gets a group by ID
func (s *MessageGroupService) GetGroupByID(id uuid.UUID) (*models.MessageGroup, error) {
	return s.groupRepo.GetByID(id)
}

// GetGroupBySlug gets a group by slug
func (s *MessageGroupService) GetGroupBySlug(slug string) (*models.MessageGroup, error) {
	return s.groupRepo.GetBySlug(slug)
}

// GetGroupsByUserID gets all groups for a user
func (s *MessageGroupService) GetGroupsByUserID(userID uuid.UUID) ([]models.MessageGroup, error) {
	return s.groupRepo.GetByUserID(userID)
}

// GetActiveGroupsByUserID gets all active groups for a user
func (s *MessageGroupService) GetActiveGroupsByUserID(userID uuid.UUID) ([]models.MessageGroup, error) {
	return s.groupRepo.GetActiveByUserID(userID)
}

// UpdateGroup updates a group
func (s *MessageGroupService) UpdateGroup(id uuid.UUID, name, description string, isPublic bool, settings map[string]interface{}) error {
	// Get group
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Convert settings to JSON
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	// Update group
	group.Name = name
	group.Description = description
	group.IsPublic = isPublic
	group.Settings = settingsJSON
	group.UpdatedAt = time.Now()

	// Save group
	return s.groupRepo.Update(group)
}

// ArchiveGroup archives a group
func (s *MessageGroupService) ArchiveGroup(id uuid.UUID) error {
	// Get group
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if already archived
	if group.IsArchived {
		return nil
	}

	// Archive group
	if err := s.groupRepo.Archive(id); err != nil {
		return err
	}

	// Decrement user's active groups count
	return s.userRepo.DecrementActiveGroups(group.UserID)
}

// UnarchiveGroup unarchives a group
func (s *MessageGroupService) UnarchiveGroup(id uuid.UUID) error {
	// Get group
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if already active
	if !group.IsArchived {
		return nil
	}

	// Check if user can activate another group
	canCreate, err := s.canUserCreateGroup(group.UserID)
	if err != nil {
		return err
	}
	if !canCreate {
		return errors.New("user has reached the maximum number of active groups")
	}

	// Unarchive group
	if err := s.groupRepo.Unarchive(id); err != nil {
		return err
	}

	// Increment user's active groups count
	return s.userRepo.IncrementActiveGroups(group.UserID)
}

// DeleteGroup deletes a group
func (s *MessageGroupService) DeleteGroup(id uuid.UUID) error {
	// Get group
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if active
	if !group.IsArchived {
		// Decrement user's active groups count
		if err := s.userRepo.DecrementActiveGroups(group.UserID); err != nil {
			return err
		}
	}

	// Delete group
	return s.groupRepo.Delete(id)
}

// UpdateGroupSettings updates a group's settings
func (s *MessageGroupService) UpdateGroupSettings(id uuid.UUID, settings map[string]interface{}) error {
	// Get group
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Convert settings to JSON
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	// Update group
	group.Settings = settingsJSON
	group.UpdatedAt = time.Now()

	// Save group
	return s.groupRepo.Update(group)
}

// IsUserOwner checks if a user is the owner of a group
func (s *MessageGroupService) IsUserOwner(groupID, userID uuid.UUID) (bool, error) {
	// Get group
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return false, err
	}

	return group.UserID == userID, nil
}

// generateSlug generates a unique slug for a group
func (s *MessageGroupService) generateSlug(name string) string {
	// Convert name to lowercase and replace spaces with hyphens
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)

	// Ensure slug is unique
	baseSlug := slug
	counter := 1
	for {
		// Check if slug is available
		available, err := s.groupRepo.IsSlugAvailable(slug)
		if err != nil || available {
			break
		}

		// Append counter to slug
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return slug
}

// canUserCreateGroup checks if a user can create a new group
func (s *MessageGroupService) canUserCreateGroup(userID uuid.UUID) (bool, error) {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	// Premium users can create unlimited groups
	if user.IsPremium() {
		return true, nil
	}

	// Get active groups count
	count, err := s.groupRepo.CountActiveByUserID(userID)
	if err != nil {
		return false, err
	}

	// Check if user has reached the group limit
	return int(count) < user.GetGroupLimit(), nil
}
