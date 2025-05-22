package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orinicee/finanzas/internal/domain"
)

type TransactionHandler struct {
	transactionUseCase domain.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase domain.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

// RegisterRoutes registra las rutas de transacciones
func (h *TransactionHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware, transactionAuthMiddleware gin.HandlerFunc) {
	transactions := router.Group("/transactions")
	transactions.Use(authMiddleware)
	{
		transactions.POST("", h.CreateTransaction)
		transactions.GET("", h.GetTransactionsByUserID)
		transactions.GET("/date-range", h.GetTransactionsByDateRange)
		transactions.GET("/category/:category", h.GetTransactionsByCategory)
		transactions.GET("/type/:type", h.GetTransactionsByType)
		transactions.GET("/recurring", h.GetRecurringTransactions)
		transactions.GET("/tax-deductible", h.GetTaxDeductibleTransactions)

		// Rutas que requieren ID específico
		transactionsWithID := transactions.Group("/:id")
		transactionsWithID.Use(transactionAuthMiddleware)
		{
			transactionsWithID.GET("", h.GetTransactionByID)
			transactionsWithID.PUT("", h.UpdateTransaction)
			transactionsWithID.DELETE("", h.DeleteTransaction)
		}
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction domain.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Asignar el ID del usuario autenticado
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)
	transaction.UserID = userID

	if err := h.transactionUseCase.CreateTransaction(c.Request.Context(), &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	transaction, err := h.transactionUseCase.GetTransactionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) GetTransactionsByUserID(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)

	transactions, err := h.transactionUseCase.GetTransactionsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionsByDateRange(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)

	startDate, err := time.Parse(time.RFC3339, c.Query("start_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha inicial inválida"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha final inválida"})
		return
	}

	transactions, err := h.transactionUseCase.GetTransactionsByDateRange(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var transaction domain.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction.ID = id
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)
	transaction.UserID = userID

	if err := h.transactionUseCase.UpdateTransaction(c.Request.Context(), &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.transactionUseCase.DeleteTransaction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TransactionHandler) GetTransactionsByCategory(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)
	category := domain.TransactionCategory(c.Param("category"))

	transactions, err := h.transactionUseCase.GetTransactionsByCategory(c.Request.Context(), userID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionsByType(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)
	transactionType := domain.TransactionType(c.Param("type"))

	transactions, err := h.transactionUseCase.GetTransactionsByType(c.Request.Context(), userID, transactionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetRecurringTransactions(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)

	transactions, err := h.transactionUseCase.GetRecurringTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTaxDeductibleTransactions(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)

	transactions, err := h.transactionUseCase.GetTaxDeductibleTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
