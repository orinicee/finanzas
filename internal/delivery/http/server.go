package http

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/orinicee/finanzas/internal/application/auth"
	"github.com/orinicee/finanzas/internal/application/report"
	"github.com/orinicee/finanzas/internal/application/transaction"
	"github.com/orinicee/finanzas/internal/application/user"
	"github.com/orinicee/finanzas/internal/delivery/http/middleware"
	"github.com/orinicee/finanzas/internal/domain"
	"github.com/orinicee/finanzas/internal/infrastructure/database"
	"github.com/orinicee/finanzas/internal/infrastructure/repository/postgres"
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

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server := &Server{
		router: router,
		config: cfg,
	}

	// Inicializar casos de uso
	userRepo := database.NewUserRepository(db.GetDB())
	userUseCase := user.NewUserUseCase(userRepo)
	authUseCase := auth.NewAuthUseCase(userRepo, cfg.JWTKey)
	transactionRepo := postgres.NewTransactionRepository(db.GetDB())
	transactionUseCase := transaction.NewTransactionUseCase(transactionRepo)
	reportUseCase := report.NewReportUseCase(transactionRepo)

	// Inicializar handlers
	userHandler := NewUserHandler(userUseCase, userRepo)
	authHandler := NewAuthHandler(authUseCase)
	transactionHandler := NewTransactionHandler(transactionUseCase)
	reportHandler := NewReportHandler(reportUseCase)

	// Configurar rutas
	server.setupRoutes(userHandler, authHandler, transactionHandler, reportHandler, authUseCase, transactionUseCase)

	return server
}

// setupRoutes configura las rutas del servidor
func (s *Server) setupRoutes(
	userHandler *UserHandler,
	authHandler *AuthHandler,
	transactionHandler *TransactionHandler,
	reportHandler *ReportHandler,
	authUseCase domain.AuthUseCase,
	transactionUseCase domain.TransactionUseCase,
) {
	// Health check (fuera del grupo v1 y sin autenticación)
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Grupo de rutas API v1
	api := s.router.Group("/api/v1")
	{
		// Rutas de autenticación
		authHandler.RegisterRoutes(api)

		// Rutas de usuarios
		userHandler.RegisterRoutes(api, middleware.AuthMiddleware(authUseCase))

		// Rutas de transacciones
		transactionHandler.RegisterRoutes(api,
			middleware.AuthMiddleware(authUseCase),
			middleware.TransactionAuthMiddleware(transactionUseCase),
		)

		// Rutas de reportes
		reportHandler.RegisterRoutes(api, middleware.AuthMiddleware(authUseCase))
	}
}

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	log.Printf("Servidor iniciado en http://%s", addr)
	return s.router.Run(addr)
}
