package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ralfferreira/papo-reto/internal/auth"
	"github.com/ralfferreira/papo-reto/internal/config"
	"github.com/ralfferreira/papo-reto/internal/handlers"
	"github.com/ralfferreira/papo-reto/internal/middleware"
	"github.com/ralfferreira/papo-reto/internal/repository"
	"github.com/ralfferreira/papo-reto/internal/services"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
	db     *repository.Database
}

// NewServer creates a new server
func NewServer(cfg *config.Config, db *repository.Database) *Server {
	// Create router
	router := gin.Default()

	// Adicionar middleware CORS
	router.Use(middleware.CORSMiddleware())

	// Create JWT service
	jwtService := auth.NewJWTService(cfg)

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Create repositories
	userRepo := repository.NewUserRepository(db.DB)
	groupRepo := repository.NewMessageGroupRepository(db.DB)
	messageRepo := repository.NewMessageRepository(db.DB)
	sharedAccessRepo := repository.NewSharedAccessRepository(db.DB)

	// Create services
	userService := services.NewUserService(userRepo, jwtService)
	groupService := services.NewMessageGroupService(groupRepo, userRepo)

	// Create handlers
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)
	groupHandler := handlers.NewGroupHandler(groupService)

	// Public routes
	router.POST("/api/v1/auth/register", authHandler.Register)
	router.POST("/api/v1/auth/login", authHandler.Login)
	router.POST("/api/v1/auth/refresh", authHandler.RefreshToken)

	// Public message sending endpoint
	router.POST("/api/v1/public/send/:slug", handlers.SendAnonymousMessage(messageRepo, groupRepo, userRepo))

	// Protected routes
	api := router.Group("/api/v1")
	api.Use(authMiddleware.RequireAuth())
	{
		// User routes
		api.GET("/user/profile", userHandler.GetProfile)
		api.PUT("/user/profile", userHandler.UpdateProfile)
		api.PUT("/user/password", userHandler.UpdatePassword)
		api.PUT("/user/notifications", userHandler.UpdateNotifications)

		// Group routes
		api.GET("/groups", groupHandler.GetGroups)
		api.POST("/groups", groupHandler.CreateGroup)
		api.GET("/groups/:id", groupHandler.GetGroup)
		api.PUT("/groups/:id", groupHandler.UpdateGroup)
		api.DELETE("/groups/:id", groupHandler.ArchiveGroup)
		api.POST("/groups/:id/unarchive", groupHandler.UnarchiveGroup)

		// Message routes
		api.GET("/groups/:id/messages", handlers.GetMessages(messageRepo))
		api.PUT("/messages/:id", handlers.UpdateMessage(messageRepo))
		api.DELETE("/messages/:id", handlers.DeleteMessage(messageRepo))

		// Shared access routes
		api.POST("/groups/:id/share", handlers.CreateSharedAccess(sharedAccessRepo, groupRepo))
		api.GET("/groups/:id/shared", handlers.GetSharedAccess(sharedAccessRepo))
		api.DELETE("/groups/:id/share/:shareId", handlers.RevokeSharedAccess(sharedAccessRepo))
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{
		config: cfg,
		router: router,
		server: server,
		db:     db,
	}
}

// Start starts the server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}
