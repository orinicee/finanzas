package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/orinicee/finanzas/internal/domain"
)

// TransactionRepository define la interfaz para el repositorio de transacciones
type TransactionRepository interface {
	// Create crea una nueva transacción
	Create(ctx context.Context, transaction *domain.Transaction) error

	// GetByID obtiene una transacción por su ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)

	// GetByUserID obtiene todas las transacciones de un usuario
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error)

	// GetByDateRange obtiene las transacciones de un usuario en un rango de fechas
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*domain.Transaction, error)

	// Update actualiza una transacción existente
	Update(ctx context.Context, transaction *domain.Transaction) error

	// Delete elimina una transacción
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByCategory obtiene las transacciones de un usuario por categoría
	GetByCategory(ctx context.Context, userID uuid.UUID, category domain.TransactionCategory) ([]*domain.Transaction, error)

	// GetByType obtiene las transacciones de un usuario por tipo (ingreso/egreso)
	GetByType(ctx context.Context, userID uuid.UUID, transactionType domain.TransactionType) ([]*domain.Transaction, error)

	// GetRecurring obtiene las transacciones recurrentes de un usuario
	GetRecurring(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error)

	// GetTaxDeductible obtiene las transacciones deducibles de impuestos de un usuario
	GetTaxDeductible(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error)
}
