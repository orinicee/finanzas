package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/orinicee/finanzas/internal/domain"
	"github.com/orinicee/finanzas/internal/infrastructure/database/migrations"
	"github.com/orinicee/finanzas/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresRepository implementa la interfaz UserRepository para PostgreSQL
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository crea una nueva instancia del repositorio PostgreSQL
func NewPostgresRepository(cfg *config.Config) (domain.UserRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error conectando a la base de datos: %v", err)
	}

	// Ejecutar migraciones
	if err := migrations.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("error ejecutando migraciones: %v", err)
	}

	log.Println("Conexión a la base de datos establecida exitosamente")
	return &PostgresRepository{db: db}, nil
}

// GetDB implementa el método de la interfaz Repository
func (r *PostgresRepository) GetDB() *gorm.DB {
	return r.db
}

// Close implementa el método de la interfaz Repository
func (r *PostgresRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Create implementa el método de la interfaz UserRepository
func (r *PostgresRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByID implementa el método de la interfaz UserRepository
func (r *PostgresRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail implementa el método de la interfaz UserRepository
func (r *PostgresRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update implementa el método de la interfaz UserRepository
func (r *PostgresRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete implementa el método de la interfaz UserRepository
func (r *PostgresRepository) Delete(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

// List implementa el método de la interfaz UserRepository
func (r *PostgresRepository) List() ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// WithTransaction implementa el método de la interfaz Transaction
func (r *PostgresRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// FindBySocialID implementa el método de la interfaz UserRepository
// func (r *PostgresRepository) FindBySocialID(provider domain.AuthProvider, socialID string) (*domain.User, error) {
// 	var user domain.User
// 	if err := r.db.First(&user, "provider = ? AND social_id = ?", provider, socialID).Error; err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// SaveRefreshToken implementa el método de la interfaz UserRepository
func (r *PostgresRepository) SaveRefreshToken(userID string, token string, expiresAt time.Time) error {
	// Definir la estructura de la tabla de refresh tokens
	type RefreshToken struct {
		Token     string    `gorm:"primaryKey;column:token"`
		UserID    string    `gorm:"index;column:user_id"`
		ExpiresAt time.Time `gorm:"column:expires_at"`
	}

	// Auto-migrar la tabla de refresh tokens
	if err := r.db.AutoMigrate(&RefreshToken{}); err != nil {
		return err
	}

	// Crear el nuevo refresh token
	refreshToken := RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	return r.db.Create(&refreshToken).Error
}

// GetRefreshToken implementa el método de la interfaz UserRepository
func (r *PostgresRepository) GetRefreshToken(token string) (string, error) {
	type RefreshToken struct {
		Token     string    `gorm:"primaryKey;column:token"`
		UserID    string    `gorm:"index;column:user_id"`
		ExpiresAt time.Time `gorm:"column:expires_at"`
	}

	var refreshToken RefreshToken
	if err := r.db.First(&refreshToken, "token = ? AND expires_at > ?", token, time.Now()).Error; err != nil {
		return "", err
	}
	return refreshToken.UserID, nil
}

// DeleteRefreshToken implementa el método de la interfaz UserRepository
func (r *PostgresRepository) DeleteRefreshToken(token string) error {
	type RefreshToken struct {
		Token     string    `gorm:"primaryKey;column:token"`
		UserID    string    `gorm:"index;column:user_id"`
		ExpiresAt time.Time `gorm:"column:expires_at"`
	}

	return r.db.Delete(&RefreshToken{}, "token = ?", token).Error
}
