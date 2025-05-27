package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/orinicee/finanzas/internal/application/auth"
	"github.com/orinicee/finanzas/internal/domain"
)

// AuthHandler maneja las solicitudes HTTP relacionadas con la autenticación
type AuthHandler struct {
	authUseCase domain.AuthUseCase
}

// NewAuthHandler crea una nueva instancia del manejador de autenticación
func NewAuthHandler(authUseCase domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
	}
}

// Register maneja el registro de nuevos usuarios
func (h *AuthHandler) Register(c *gin.Context) {
	var credentials domain.AuthCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el cuerpo de la solicitud"})
		return
	}

	fmt.Printf("📥 Datos recibidos en el backend (Register): %+v\n", credentials)

	token, err := h.authUseCase.Register(&credentials)
	if err != nil {
		if err == domain.ErrEmailAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El email ya está registrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el usuario"})
		fmt.Printf("❌ Error en authUseCase.Register: %v\n", err)
		return
	}

	// Obtener el usuario recién creado
	user, err := h.authUseCase.ValidateToken(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del usuario"})
		return
	}

	response := gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"token_type":    token.TokenType,
		"expires_in":    token.ExpiresIn,
	}

	fmt.Printf("📤 Respuesta enviada al front (Register): %+v\n", response)
	c.JSON(http.StatusCreated, response)
}

// Login maneja el inicio de sesión de usuarios
func (h *AuthHandler) Login(c *gin.Context) {
	var credentials domain.AuthCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el cuerpo de la solicitud"})
		return
	}

	token, err := h.authUseCase.Login(&credentials)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al iniciar sesión"})
		return
	}

	user, err := h.authUseCase.ValidateToken(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del usuario"})
		return
	}

	response := gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"token_type":    token.TokenType,
		"expires_in":    token.ExpiresIn,
	}

	fmt.Printf("📤 Respuesta enviada al front (Login): %+v\n", response)
	c.JSON(http.StatusCreated, response)
}

// RefreshToken maneja la renovación de tokens
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token de actualización requerido"})
		return
	}

	token, err := h.authUseCase.RefreshToken(request.RefreshToken)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrTokenExpired {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al renovar el token"})
		return
	}

	response := gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"token_type":    token.TokenType,
		"expires_in":    token.ExpiresIn,
	}

	c.JSON(http.StatusOK, response)
}

// Logout maneja el cierre de sesión
func (h *AuthHandler) Logout(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token de actualización requerido"})
		return
	}

	if err := h.authUseCase.Logout(request.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cerrar sesión"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}
