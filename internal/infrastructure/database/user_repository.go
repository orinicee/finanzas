package database

import (
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
