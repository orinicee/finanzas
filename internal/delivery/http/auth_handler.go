package http

import (
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

func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/social", h.SocialAuth)
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

	token, err := h.authUseCase.Register(&credentials)
	if err != nil {
		if err == domain.ErrEmailAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El email ya está registrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el usuario"})
		return
	}

	c.JSON(http.StatusCreated, token)
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

	c.JSON(http.StatusOK, token)
}

// SocialAuth maneja la autenticación con proveedores sociales
func (h *AuthHandler) SocialAuth(c *gin.Context) {
	var credentials domain.SocialAuthCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el cuerpo de la solicitud"})
		return
	}

	token, err := h.authUseCase.SocialAuth(&credentials)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error en la autenticación social"})
		return
	}

	c.JSON(http.StatusOK, token)
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

	c.JSON(http.StatusOK, token)
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
