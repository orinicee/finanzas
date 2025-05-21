package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/orinicee/finanzas/internal/application/user"
	"github.com/orinicee/finanzas/internal/domain"
)

// UserHandler maneja las solicitudes HTTP relacionadas con usuarios
type UserHandler struct {
	userUseCase user.UserUseCase
	userRepo    domain.UserRepository
}

// NewUserHandler crea una nueva instancia del manejador de usuarios
func NewUserHandler(userUseCase user.UserUseCase, userRepo domain.UserRepository) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		userRepo:    userRepo,
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	users := router.Group("/api/v1/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.GET("", h.ListUsers)
	}
}

// CreateUser maneja la creación de un nuevo usuario
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el cuerpo de la solicitud"})
		return
	}

	if err := h.userUseCase.CreateUser(&user); err != nil {
		if err == domain.ErrEmailAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El email ya está registrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado exitosamente"})
}

// GetUser maneja la obtención de un usuario por ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario requerido"})
		return
	}

	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el usuario"})
		return
	}

	// Devolver solo los campos necesarios
	response := gin.H{
		"id":              user.ID,
		"full_name":       user.FullName,
		"email":           user.Email,
		"document_type":   user.DocumentType,
		"document_number": user.DocumentNumber,
		"tax_regime":      user.TaxRegime,
		"person_type":     user.PersonType,
		"city":            user.City,
		"department":      user.Department,
		"address":         user.Address,
		"phone":           user.Phone,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser maneja la actualización de un usuario
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el cuerpo de la solicitud"})
		return
	}

	user.ID = id
	if err := h.userUseCase.UpdateUser(&user); err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado exitosamente"})
}

// DeleteUser maneja la eliminación de un usuario
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userUseCase.DeleteUser(id); err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}

// ListUsers maneja la obtención de todos los usuarios
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userUseCase.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la lista de usuarios"})
		return
	}

	// Devolver solo los campos necesarios para cada usuario
	response := make([]gin.H, 0)
	for _, user := range users {
		response = append(response, gin.H{
			"id":              user.ID,
			"full_name":       user.FullName,
			"email":           user.Email,
			"document_type":   user.DocumentType,
			"document_number": user.DocumentNumber,
			"tax_regime":      user.TaxRegime,
			"person_type":     user.PersonType,
			"city":            user.City,
			"department":      user.Department,
			"address":         user.Address,
			"phone":           user.Phone,
		})
	}

	c.JSON(http.StatusOK, response)
}
