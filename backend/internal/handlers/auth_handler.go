package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ralfferreira/papo-reto/internal/services"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userService *services.UserService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	// Parse request
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Register user
	user, err := h.userService.RegisterUser(req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return user without password
	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"isVerified": user.IsVerified,
		"plan":       user.Plan,
		"createdAt":  user.CreatedAt,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	// Parse request
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Login user
	token, err := h.userService.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Parse request
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Refresh token
	newToken, err := h.userService.RefreshToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return new token
	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
