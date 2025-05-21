package http

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/orinicee/finanzas/internal/application/user"
	"github.com/orinicee/finanzas/internal/infrastructure/database"
	"github.com/orinicee/finanzas/pkg/config"
)

// Server representa el servidor HTTP
type Server struct {
	router *gin.Engine
	config *config.Config
}

// NewServer crea una nueva instancia del servidor
func NewServer(cfg *config.Config, db *database.PostgresRepository) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server := &Server{
		router: router,
		config: cfg,
	}

	// Inicializar casos de uso
	userRepo := database.NewUserRepository(db.GetDB())
	userUseCase := user.NewUserUseCase(userRepo)

	// Inicializar handlers
	userHandler := NewUserHandler(userUseCase, userRepo)

	// Configurar rutas
	server.setupRoutes(userHandler)

	return server
}

// setupRoutes configura las rutas del servidor
func (s *Server) setupRoutes(userHandler *UserHandler) {
	// Grupo de rutas API v1
	v1 := s.router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		// Rutas de usuarios
		userHandler.RegisterRoutes(s.router)
	}
}

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	log.Printf("Servidor iniciado en http://%s", addr)
	return s.router.Run(addr)
}
