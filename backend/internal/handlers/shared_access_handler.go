package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/models"
	"github.com/ralfferreira/papo-reto/internal/repository"
)

// CreateSharedAccess returns a handler for creating shared access to a group
func CreateSharedAccess(sharedAccessRepo *repository.SharedAccessRepository, groupRepo *repository.MessageGroupRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Get group ID from URL
		groupID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
			return
		}

		// Check if group exists and user is the owner
		group, err := groupRepo.GetByID(groupID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}

		if group.UserID != userID.(uuid.UUID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to share this group"})
			return
		}

		// Parse request
		var req struct {
			Email     string     `json:"email" binding:"required,email"`
			ExpiresAt *time.Time `json:"expiresAt"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate token
		token := uuid.New().String()

		// Create shared access
		sharedAccess := &models.SharedAccess{
			GroupID:   groupID,
			InvitedBy: userID.(uuid.UUID),
			Email:     req.Email,
			Token:     token,
			IsActive:  true,
			ExpiresAt: req.ExpiresAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Save shared access
		if err := sharedAccessRepo.Create(sharedAccess); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return shared access
		c.JSON(http.StatusCreated, gin.H{
			"id":        sharedAccess.ID,
			"groupId":   sharedAccess.GroupID,
			"email":     sharedAccess.Email,
			"token":     sharedAccess.Token,
			"isActive":  sharedAccess.IsActive,
			"expiresAt": sharedAccess.ExpiresAt,
			"createdAt": sharedAccess.CreatedAt,
		})
	}
}

// GetSharedAccess returns a handler for getting shared access for a group
func GetSharedAccess(sharedAccessRepo *repository.SharedAccessRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Get group ID from URL
		groupID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
			return
		}

		// Get shared access
		sharedAccess, err := sharedAccessRepo.GetByGroupID(groupID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Convert to response format
		var response []gin.H
		for _, access := range sharedAccess {
			response = append(response, gin.H{
				"id":        access.ID,
				"email":     access.Email,
				"token":     access.Token,
				"isActive":  access.IsActive,
				"expiresAt": access.ExpiresAt,
				"createdAt": access.CreatedAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{"sharedAccess": response})
	}
}

// RevokeSharedAccess returns a handler for revoking shared access
func RevokeSharedAccess(sharedAccessRepo *repository.SharedAccessRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Get group ID and share ID from URL
		_, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
			return
		}

		shareID, err := uuid.Parse(c.Param("shareId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share ID"})
			return
		}

		// Revoke shared access
		if err := sharedAccessRepo.Revoke(shareID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "shared access revoked successfully"})
	}
}
