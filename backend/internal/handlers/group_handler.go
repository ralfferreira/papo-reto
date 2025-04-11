package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/services"
)

// GroupHandler handles group requests
type GroupHandler struct {
	groupService *services.MessageGroupService
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(groupService *services.MessageGroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// GetGroups handles getting all groups for a user
func (h *GroupHandler) GetGroups(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check if we should include archived groups
	includeArchived := c.Query("includeArchived") == "true"

	// Get groups
	var groups []interface{}

	if includeArchived {
		// Get all groups
		allGroups, _ := h.groupService.GetGroupsByUserID(userID.(uuid.UUID))

		// Convert to response format
		for _, group := range allGroups {
			groups = append(groups, gin.H{
				"id":          group.ID,
				"name":        group.Name,
				"slug":        group.Slug,
				"description": group.Description,
				"isPublic":    group.IsPublic,
				"isArchived":  group.IsArchived,
				"createdAt":   group.CreatedAt,
			})
		}
	} else {
		// Get active groups
		activeGroups, _ := h.groupService.GetActiveGroupsByUserID(userID.(uuid.UUID))

		// Convert to response format
		for _, group := range activeGroups {
			groups = append(groups, gin.H{
				"id":          group.ID,
				"name":        group.Name,
				"slug":        group.Slug,
				"description": group.Description,
				"isPublic":    group.IsPublic,
				"isArchived":  group.IsArchived,
				"createdAt":   group.CreatedAt,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// CreateGroup handles creating a new group
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse request
	var req struct {
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		IsPublic    bool                   `json:"isPublic"`
		Settings    map[string]interface{} `json:"settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create group
	group, err := h.groupService.CreateGroup(userID.(uuid.UUID), req.Name, req.Description, req.IsPublic, req.Settings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse settings
	var settings map[string]interface{}
	if group.Settings != nil {
		if err := json.Unmarshal(group.Settings, &settings); err != nil {
			settings = make(map[string]interface{})
		}
	} else {
		settings = make(map[string]interface{})
	}

	// Return group
	c.JSON(http.StatusCreated, gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"slug":        group.Slug,
		"description": group.Description,
		"isPublic":    group.IsPublic,
		"isArchived":  group.IsArchived,
		"settings":    settings,
		"createdAt":   group.CreatedAt,
	})
}

// GetGroup handles getting a group by ID
func (h *GroupHandler) GetGroup(c *gin.Context) {
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

	// Check if user is the owner
	isOwner, err := h.groupService.IsUserOwner(groupID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to access this group"})
		return
	}

	// Get group
	group, err := h.groupService.GetGroupByID(groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Parse settings
	var settings map[string]interface{}
	if group.Settings != nil {
		if err := json.Unmarshal(group.Settings, &settings); err != nil {
			settings = make(map[string]interface{})
		}
	} else {
		settings = make(map[string]interface{})
	}

	// Return group
	c.JSON(http.StatusOK, gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"slug":        group.Slug,
		"description": group.Description,
		"isPublic":    group.IsPublic,
		"isArchived":  group.IsArchived,
		"settings":    settings,
		"createdAt":   group.CreatedAt,
	})
}

// UpdateGroup handles updating a group
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
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

	// Check if user is the owner
	isOwner, err := h.groupService.IsUserOwner(groupID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to update this group"})
		return
	}

	// Parse request
	var req struct {
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		IsPublic    bool                   `json:"isPublic"`
		Settings    map[string]interface{} `json:"settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update group
	if err := h.groupService.UpdateGroup(groupID, req.Name, req.Description, req.IsPublic, req.Settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "group updated successfully"})
}

// ArchiveGroup handles archiving a group
func (h *GroupHandler) ArchiveGroup(c *gin.Context) {
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

	// Check if user is the owner
	isOwner, err := h.groupService.IsUserOwner(groupID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to archive this group"})
		return
	}

	// Archive group
	if err := h.groupService.ArchiveGroup(groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "group archived successfully"})
}

// UnarchiveGroup handles unarchiving a group
func (h *GroupHandler) UnarchiveGroup(c *gin.Context) {
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

	// Check if user is the owner
	isOwner, err := h.groupService.IsUserOwner(groupID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to unarchive this group"})
		return
	}

	// Unarchive group
	if err := h.groupService.UnarchiveGroup(groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "group unarchived successfully"})
}
