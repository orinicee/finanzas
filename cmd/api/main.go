package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/orinicee/finanzas/internal/application/auth"
	"github.com/orinicee/finanzas/internal/application/user"
	"github.com/orinicee/finanzas/internal/delivery/http"
	"github.com/orinicee/finanzas/internal/delivery/http/middleware"
	"github.com/orinicee/finanzas/internal/infrastructure/database"
	"github.com/orinicee/finanzas/pkg/config"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Inicializar repositorio de base de datos
	userRepo, err := database.NewPostgresRepository(cfg)
	if err != nil {
		log.Fatalf("Error inicializando repositorio: %v", err)
	}
	defer userRepo.Close()

	// Inicializar caso de uso
	userUseCase := user.NewUserUseCase(userRepo)
	authUseCase := auth.NewAuthUseCase(userRepo, cfg.JWTKey)

	// Inicializar Gin
	router := gin.Default()

	// Configurar CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Endpoint de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Grupo de rutas API v1
	api := router.Group("/api/v1")

	// Inicializar manejadores
	userHandler := http.NewUserHandler(userUseCase, userRepo)
	authHandler := http.NewAuthHandler(authUseCase)

	// Registrar rutas
	authHandler.RegisterRoutes(api)
	userHandler.RegisterRoutes(api, middleware.AuthMiddleware(authUseCase))

	// Iniciar servidor
	serverAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Printf("Servidor iniciado en http://%s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
