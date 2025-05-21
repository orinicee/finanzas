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

	// Inicializar Gin
	router := gin.Default()

	// Inicializar manejador de usuarios
	userHandler := http.NewUserHandler(userUseCase, userRepo)

	authUseCase := auth.NewAuthUseCase(userRepo, cfg.JWTKey)

	// Registrar rutas
	userHandler.RegisterRoutes(router, middleware.AuthMiddleware(authUseCase))

	// Iniciar servidor
	serverAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Printf("Servidor iniciado en http://%s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
