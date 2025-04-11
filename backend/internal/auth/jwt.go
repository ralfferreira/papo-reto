package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ralfferreira/papo-reto/internal/config"
	"github.com/ralfferreira/papo-reto/internal/models"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// JWTService handles JWT operations
type JWTService struct {
	config *config.Config
}

// NewJWTService creates a new JWT service
func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		config: cfg,
	}
}

// GenerateToken generates a new JWT token for a user
func (s *JWTService) GenerateToken(user *models.User) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(time.Duration(s.config.JWT.ExpiryMinutes) * time.Minute)

	// Create claims
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "papo-reto-api",
			Subject:   user.ID.String(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken refreshes a JWT token
func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	// Validate existing token
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Create a new token with the same claims but extended expiration
	expirationTime := time.Now().Add(time.Duration(s.config.JWT.ExpiryMinutes) * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.NotBefore = jwt.NewNumericDate(time.Now())

	// Create new token with refreshed claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	newTokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}
