package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/orinicee/finanzas/internal/domain"
	"github.com/orinicee/finanzas/internal/domain/repository"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository crea una nueva instancia del repositorio de transacciones
func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

func (r *transactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Transaction{}, "id = ?", id).Error
}

func (r *transactionRepository) GetByCategory(ctx context.Context, userID uuid.UUID, category domain.TransactionCategory) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND category = ?", userID, category).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetByType(ctx context.Context, userID uuid.UUID, transactionType domain.TransactionType) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, transactionType).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetRecurring(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_recurring = true", userID).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetTaxDeductible(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_tax_deductible = true", userID).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
