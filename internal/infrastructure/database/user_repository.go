package database

import (
	"time"

	"github.com/orinicee/finanzas/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	// Auto-migrar la tabla de usuarios
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		panic("Error al migrar la tabla de usuarios: " + err.Error())
	}

	return &userRepository{db: db}
}

// GetDB implementa el método de la interfaz Repository
func (r *userRepository) GetDB() *gorm.DB {
	return r.db
}

// Close implementa el método de la interfaz Repository
func (r *userRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindBySocialID(provider domain.AuthProvider, socialID string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "provider = ? AND social_id = ?", provider, socialID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepository) List() ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Métodos para manejo de refresh tokens
type refreshToken struct {
	Token     string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	ExpiresAt time.Time
}

func (r *userRepository) SaveRefreshToken(userID string, token string, expiresAt time.Time) error {
	// Auto-migrar la tabla de refresh tokens
	if err := r.db.AutoMigrate(&refreshToken{}); err != nil {
		return err
	}

	return r.db.Create(&refreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}).Error
}

func (r *userRepository) GetRefreshToken(token string) (string, error) {
	var rt refreshToken
	err := r.db.First(&rt, "token = ? AND expires_at > ?", token, time.Now()).Error
	if err != nil {
		return "", err
	}
	return rt.UserID, nil
}

func (r *userRepository) DeleteRefreshToken(token string) error {
	return r.db.Delete(&refreshToken{}, "token = ?", token).Error
}
