package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orinicee/finanzas/internal/domain"
)

// AuthMiddleware es un middleware que verifica la autenticación del usuario
func AuthMiddleware(authUseCase domain.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no se proporcionó token de autorización"})
			c.Abort()
			return
		}

		// Extraer el token del header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			c.Abort()
			return
		}

		// Validar el token
		user, err := authUseCase.ValidateToken(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			c.Abort()
			return
		}

		// Almacenar el usuario en el contexto
		c.Set("user", user)
		c.Next()
	}
}

// TransactionAuthMiddleware es un middleware que verifica que la transacción pertenece al usuario autenticado
func TransactionAuthMiddleware(transactionUseCase domain.TransactionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el usuario del contexto
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
			c.Abort()
			return
		}

		// Obtener el ID de la transacción de los parámetros
		transactionID := c.Param("id")
		if transactionID == "" {
			c.Next()
			return
		}

		// Convertir el ID a UUID
		id, err := uuid.Parse(transactionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de transacción inválido"})
			c.Abort()
			return
		}

		// Obtener la transacción
		transaction, err := transactionUseCase.GetTransactionByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "transacción no encontrada"})
			c.Abort()
			return
		}

		// Verificar que la transacción pertenece al usuario
		userID, err := uuid.Parse(user.(*domain.User).ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al procesar el ID del usuario"})
			c.Abort()
			return
		}
		if transaction.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "no autorizado para acceder a esta transacción"})
			c.Abort()
			return
		}

		c.Next()
	}
}
