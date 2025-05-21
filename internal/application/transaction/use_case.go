package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/orinicee/finanzas/internal/domain"
)

var (
	ErrTransactionNotFound = errors.New("transacción no encontrada")
	ErrInvalidAmount       = errors.New("monto inválido")
	ErrInvalidDate         = errors.New("fecha inválida")
	ErrUnauthorized        = errors.New("no autorizado para acceder a esta transacción")
)

type UseCase struct {
	transactionRepo domain.TransactionRepository
}

func NewTransactionUseCase(transactionRepo domain.TransactionRepository) domain.TransactionUseCase {
	return &UseCase{
		transactionRepo: transactionRepo,
	}
}

func (uc *UseCase) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	// Validaciones básicas
	if transaction.Amount <= 0 {
		return ErrInvalidAmount
	}

	if transaction.Date.IsZero() {
		transaction.Date = time.Now()
	}

	return uc.transactionRepo.Create(ctx, transaction)
}

func (uc *UseCase) GetTransactionByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	transaction, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTransactionNotFound
	}
	return transaction, nil
}

func (uc *UseCase) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetByUserID(ctx, userID)
}

func (uc *UseCase) GetTransactionsByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	if startDate.After(endDate) {
		return nil, ErrInvalidDate
	}
	return uc.transactionRepo.GetByDateRange(ctx, userID, startDate, endDate)
}

func (uc *UseCase) UpdateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	// Verificar que la transacción existe
	existing, err := uc.transactionRepo.GetByID(ctx, transaction.ID)
	if err != nil {
		return ErrTransactionNotFound
	}

	// Validar que el usuario es el propietario
	if existing.UserID != transaction.UserID {
		return ErrUnauthorized
	}

	// Validaciones básicas
	if transaction.Amount <= 0 {
		return ErrInvalidAmount
	}

	return uc.transactionRepo.Update(ctx, transaction)
}

func (uc *UseCase) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	// Verificar que la transacción existe
	_, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return ErrTransactionNotFound
	}

	return uc.transactionRepo.Delete(ctx, id)
}

func (uc *UseCase) GetTransactionsByCategory(ctx context.Context, userID uuid.UUID, category domain.TransactionCategory) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetByCategory(ctx, userID, category)
}

func (uc *UseCase) GetTransactionsByType(ctx context.Context, userID uuid.UUID, transactionType domain.TransactionType) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetByType(ctx, userID, transactionType)
}

func (uc *UseCase) GetRecurringTransactions(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetRecurring(ctx, userID)
}

func (uc *UseCase) GetTaxDeductibleTransactions(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetTaxDeductible(ctx, userID)
}
