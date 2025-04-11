package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/auth"
	"github.com/ralfferreira/papo-reto/internal/models"
	"github.com/ralfferreira/papo-reto/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for users
type UserService struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository, jwtService *auth.JWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(email, password, name string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        email,
		Password:     string(hashedPassword),
		Name:         name,
		IsVerified:   false,
		Plan:         "free",
		MessageCount: 0,
		ActiveGroups: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser authenticates a user and returns a JWT token
func (s *UserService) LoginUser(email, password string) (string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// UpdateUserProfile updates a user's profile
func (s *UserService) UpdateUserProfile(id uuid.UUID, name, avatarURL string) error {
	// Get user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Update user
	user.Name = name
	user.AvatarURL = avatarURL
	user.UpdatedAt = time.Now()

	// Save user
	return s.userRepo.Update(user)
}

// UpdateUserPassword updates a user's password
func (s *UserService) UpdateUserPassword(id uuid.UUID, currentPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	// Save user
	return s.userRepo.Update(user)
}

// UpdateUserNotificationSettings updates a user's notification settings
func (s *UserService) UpdateUserNotificationSettings(id uuid.UUID, notifySettings []byte) error {
	// Get user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Update user
	user.NotifySettings = notifySettings
	user.UpdatedAt = time.Now()

	// Save user
	return s.userRepo.Update(user)
}

// UpdateUserPlan updates a user's plan
func (s *UserService) UpdateUserPlan(id uuid.UUID, plan string) error {
	if plan != "free" && plan != "premium" {
		return errors.New("invalid plan")
	}
	return s.userRepo.UpdatePlan(id, plan)
}

// VerifyUser marks a user as verified
func (s *UserService) VerifyUser(id uuid.UUID) error {
	// Get user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Update user
	user.IsVerified = true
	user.UpdatedAt = time.Now()

	// Save user
	return s.userRepo.Update(user)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}

// RefreshToken refreshes a JWT token
func (s *UserService) RefreshToken(token string) (string, error) {
	return s.jwtService.RefreshToken(token)
}

// CanCreateGroup checks if a user can create a new group
func (s *UserService) CanCreateGroup(id uuid.UUID) (bool, error) {
	// Get user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return false, err
	}

	// Premium users can create unlimited groups
	if user.IsPremium() {
		return true, nil
	}

	// Check if user has reached the group limit
	if user.ActiveGroups >= user.GetGroupLimit() {
		return false, nil
	}

	return true, nil
}

// CanSendMessage checks if a user can send a new message
func (s *UserService) CanSendMessage(id uuid.UUID) (bool, error) {
	// Get user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return false, err
	}

	// Premium users can send unlimited messages
	if user.IsPremium() {
		return true, nil
	}

	// Check if user has reached the message limit
	if user.MessageCount >= user.GetMessageLimit() {
		return false, nil
	}

	return true, nil
}
