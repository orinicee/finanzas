package domain

import (
	"context"

	"gorm.io/gorm"
)

// Repository define la interfaz base para el repositorio
type Repository interface {
	// GetDB retorna la conexión a la base de datos
	GetDB() *gorm.DB
	// Close cierra la conexión a la base de datos
	Close() error
}

// UserRepository define la interfaz para las operaciones de usuarios
type UserRepository interface {
	Repository
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List() ([]*User, error)
}

// Transaction define la interfaz para manejar transacciones
type Transaction interface {
	// WithTransaction ejecuta una función dentro de una transacción
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}
