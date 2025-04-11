package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/models"
	"github.com/ralfferreira/papo-reto/internal/repository"
)

// GetMessages returns a handler for getting messages in a group
func GetMessages(messageRepo *repository.MessageRepository) gin.HandlerFunc {
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

		// Get pagination parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

		// Validate pagination parameters
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}

		// Get messages
		messages, err := messageRepo.GetByGroupIDPaginated(groupID, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Convert to response format
		var response []gin.H
		for _, message := range messages {
			response = append(response, gin.H{
				"id":         message.ID,
				"content":    message.Content,
				"isRead":     message.IsRead,
				"isFavorite": message.IsFavorite,
				"isRevealed": message.IsRevealed,
				"senderID":   message.SenderID,
				"createdAt":  message.CreatedAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{"messages": response})
	}
}

// UpdateMessage returns a handler for updating a message
func UpdateMessage(messageRepo *repository.MessageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Get message ID from URL
		messageID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
			return
		}

		// Parse request
		var req struct {
			IsRead     *bool `json:"isRead"`
			IsFavorite *bool `json:"isFavorite"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get message
		message, err := messageRepo.GetByID(messageID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Update message
		if req.IsRead != nil && *req.IsRead != message.IsRead {
			message.IsRead = *req.IsRead
		}

		if req.IsFavorite != nil && *req.IsFavorite != message.IsFavorite {
			message.IsFavorite = *req.IsFavorite
		}

		// Save message
		if err := messageRepo.Update(message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "message updated successfully"})
	}
}

// DeleteMessage returns a handler for deleting a message
func DeleteMessage(messageRepo *repository.MessageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Get message ID from URL
		messageID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
			return
		}

		// Delete message
		if err := messageRepo.Delete(messageID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "message deleted successfully"})
	}
}

// SendAnonymousMessage returns a handler for sending an anonymous message
func SendAnonymousMessage(messageRepo *repository.MessageRepository, groupRepo *repository.MessageGroupRepository, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get slug from URL
		slug := c.Param("slug")
		if slug == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slug"})
			return
		}

		// Get group by slug
		group, err := groupRepo.GetBySlug(slug)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}

		// Check if group is archived
		if group.IsArchived {
			c.JSON(http.StatusForbidden, gin.H{"error": "this group is no longer accepting messages"})
			return
		}

		// Parse request
		var req struct {
			Content    string  `json:"content" binding:"required"`
			SenderID   *string `json:"senderId"`
			RevealName bool    `json:"revealName"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create message
		message := &models.Message{
			GroupID:   group.ID,
			Content:   req.Content,
			SenderIP:  c.ClientIP(),
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		// Handle identity revelation if requested
		if req.RevealName && req.SenderID != nil {
			message.IsRevealed = true
			message.SenderID = req.SenderID
		}

		// Save message
		if err := messageRepo.Create(message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Increment user's message count
		if err := userRepo.IncrementMessageCount(group.UserID); err != nil {
			// Log error but don't fail the request
			// TODO: Add proper logging
			// log.Printf("Failed to increment message count: %v", err)
		}

		c.JSON(http.StatusCreated, gin.H{"message": "message sent successfully"})
	}
}
