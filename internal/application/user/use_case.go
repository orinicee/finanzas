package user

import (
	"github.com/orinicee/finanzas/internal/domain"
)

// UseCase implementa la interfaz domain.UserUseCase
type UseCase struct {
	userRepo domain.UserRepository
}

// NewUserUseCase crea una nueva instancia del caso de uso de usuarios
func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &UseCase{
		userRepo: userRepo,
	}
}

// CreateUser crea un nuevo usuario
func (uc *UseCase) CreateUser(user *domain.User) error {
	return uc.userRepo.Create(user)
}

// GetUserByID obtiene un usuario por su ID
func (uc *UseCase) GetUserByID(id string) (*domain.User, error) {
	return uc.userRepo.FindByID(id)
}

// GetUserByEmail obtiene un usuario por su email
func (uc *UseCase) GetUserByEmail(email string) (*domain.User, error) {
	return uc.userRepo.FindByEmail(email)
}

// UpdateUser actualiza un usuario existente
func (uc *UseCase) UpdateUser(user *domain.User) error {
	return uc.userRepo.Update(user)
}

// DeleteUser elimina un usuario por su ID
func (uc *UseCase) DeleteUser(id string) error {
	return uc.userRepo.Delete(id)
}

// ListUsers obtiene la lista de todos los usuarios
func (uc *UseCase) ListUsers() ([]*domain.User, error) {
	return uc.userRepo.List()
}
