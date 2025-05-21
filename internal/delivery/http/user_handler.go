package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/orinicee/finanzas/internal/domain"
	"github.com/orinicee/finanzas/internal/domain/messages"
)

// CreateUserRequest representa la estructura de la solicitud de creación de usuario
type CreateUserRequest struct {
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	DocumentType   string `json:"document_type"`
	DocumentNumber string `json:"document_number"`
	TaxRegime      string `json:"tax_regime"`
	PersonType     string `json:"person_type"`
	City           string `json:"city"`
	Department     string `json:"department"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
}

// UpdateUserRequest representa la estructura de la solicitud de actualización de usuario
type UpdateUserRequest struct {
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	DocumentType   string `json:"document_type"`
	DocumentNumber string `json:"document_number"`
	TaxRegime      string `json:"tax_regime"`
	PersonType     string `json:"person_type"`
	City           string `json:"city"`
	Department     string `json:"department"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
}

// UserHandler maneja las solicitudes HTTP relacionadas con usuarios
type UserHandler struct {
	userUseCase domain.UserUseCase
	userRepo    domain.UserRepository
}

// NewUserHandler crea una nueva instancia del manejador de usuarios
func NewUserHandler(userUseCase domain.UserUseCase, userRepo domain.UserRepository) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		userRepo:    userRepo,
	}
}

// RegisterRoutes registra las rutas de usuarios
func (h *UserHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	users := router.Group("/api/v1/users")
	{
		// Rutas públicas
		users.POST("", h.CreateUser)

		// Rutas protegidas
		protected := users.Group("")
		protected.Use(authMiddleware)
		{
			protected.GET("/me", h.GetCurrentUser)
			protected.PUT("/me", h.UpdateCurrentUser)
			protected.DELETE("/me", h.DeleteCurrentUser)
		}
	}
}

// CreateUser maneja la creación de un nuevo usuario
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Manejar errores de validación específicos
		if err.Error() == "Key: 'CreateUserRequest.FullName' Error:Field validation for 'FullName' failed on the 'required' tag" {
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrRequiredFullName})
			return
		}
		if err.Error() == "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag" {
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrRequiredEmail})
			return
		}
		if err.Error() == "Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag" {
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrRequiredPassword})
			return
		}
		if err.Error() == "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag" {
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEmailFormat})
			return
		}
		if err.Error() == "Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag" {
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrPasswordTooShort})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrDecodeRequestBody})
		return
	}

	user := &domain.User{
		FullName:       req.FullName,
		Email:          req.Email,
		Password:       req.Password,
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		TaxRegime:      req.TaxRegime,
		PersonType:     req.PersonType,
		City:           req.City,
		Department:     req.Department,
		Address:        req.Address,
		Phone:          req.Phone,
	}

	if err := h.userUseCase.CreateUser(user); err != nil {
		switch err {
		case domain.ErrEmailAlreadyExists:
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrEmailAlreadyExists})
		case domain.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidInput})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrCreateUser})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": messages.SuccessUserCreated})
}

// GetUser maneja la obtención de un usuario por ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrRequiredUserID})
		return
	}

	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrUserNotFound})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrGetUser})
		}
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
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Manejar errores de validación específicos
		if err.Error() == "Key: 'UpdateUserRequest.FullName' Error:Field validation for 'FullName' failed on the 'required' tag" {
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrRequiredFullName})
			return
		}
		if err.Error() == "Key: 'UpdateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag" {
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrRequiredEmail})
			return
		}
		if err.Error() == "Key: 'UpdateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag" {
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEmailFormat})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrDecodeRequestBody})
		return
	}

	user := &domain.User{
		ID:             id,
		FullName:       req.FullName,
		Email:          req.Email,
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		TaxRegime:      req.TaxRegime,
		PersonType:     req.PersonType,
		City:           req.City,
		Department:     req.Department,
		Address:        req.Address,
		Phone:          req.Phone,
	}

	if err := h.userUseCase.UpdateUser(user); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrUserNotFound})
		case domain.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidInput})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrUpdateUser})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.SuccessUserUpdated})
}

// DeleteUser maneja la eliminación de un usuario
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userUseCase.DeleteUser(id); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrUserNotFound})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrDeleteUser})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.SuccessUserDeleted})
}

// ListUsers maneja la obtención de todos los usuarios
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userUseCase.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrListUsers})
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

// GetCurrentUser maneja la obtención del usuario actual
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)

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

// UpdateCurrentUser maneja la actualización del usuario actual
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	user.FullName = req.FullName
	user.Email = req.Email
	user.DocumentType = req.DocumentType
	user.DocumentNumber = req.DocumentNumber
	user.TaxRegime = req.TaxRegime
	user.PersonType = req.PersonType
	user.City = req.City
	user.Department = req.Department
	user.Address = req.Address
	user.Phone = req.Phone

	if err := h.userUseCase.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado exitosamente"})
}

// DeleteCurrentUser maneja la eliminación del usuario actual
func (h *UserHandler) DeleteCurrentUser(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)

	if err := h.userUseCase.DeleteUser(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}
