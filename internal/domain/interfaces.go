package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository define la interfaz para el repositorio de usuarios
type UserRepository interface {
	GetDB() *gorm.DB
	Close() error
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindBySocialID(provider AuthProvider, socialID string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List() ([]*User, error)
	SaveRefreshToken(userID string, token string, expiresAt time.Time) error
	GetRefreshToken(token string) (string, error)
	DeleteRefreshToken(token string) error
}

// UserUseCase define las operaciones de negocio para usuarios
type UserUseCase interface {
	CreateUser(user *User) error
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
	ListUsers() ([]*User, error)
}

// AuthUseCase define la interfaz para los casos de uso de autenticación
type AuthUseCase interface {
	Register(credentials *AuthCredentials) (*AuthToken, error)
	Login(credentials *AuthCredentials) (*AuthToken, error)
	SocialAuth(credentials *SocialAuthCredentials) (*AuthToken, error)
	RefreshToken(refreshToken string) (*AuthToken, error)
	Logout(refreshToken string) error
	ValidateToken(token string) (*User, error)
}

// TransactionRepository define la interfaz para el repositorio de transacciones
type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Transaction, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*Transaction, error)
	Update(ctx context.Context, transaction *Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCategory(ctx context.Context, userID uuid.UUID, category TransactionCategory) ([]*Transaction, error)
	GetByType(ctx context.Context, userID uuid.UUID, transactionType TransactionType) ([]*Transaction, error)
	GetRecurring(ctx context.Context, userID uuid.UUID) ([]*Transaction, error)
	GetTaxDeductible(ctx context.Context, userID uuid.UUID) ([]*Transaction, error)
}

// TransactionUseCase define las operaciones de negocio para transacciones
type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*Transaction, error)
	GetTransactionsByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
	GetTransactionsByCategory(ctx context.Context, userID uuid.UUID, category TransactionCategory) ([]*Transaction, error)
	GetTransactionsByType(ctx context.Context, userID uuid.UUID, transactionType TransactionType) ([]*Transaction, error)
	GetRecurringTransactions(ctx context.Context, userID uuid.UUID) ([]*Transaction, error)
	GetTaxDeductibleTransactions(ctx context.Context, userID uuid.UUID) ([]*Transaction, error)
}
