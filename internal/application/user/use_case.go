package user

import (
	"github.com/orinicee/finanzas/internal/domain"
)

type useCase struct {
	userRepo domain.UserRepository
}

// NewUserUseCase crea una nueva instancia del caso de uso de usuarios
func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &useCase{
		userRepo: userRepo,
	}
}

func (uc *useCase) CreateUser(user *domain.User) error {
	return uc.userRepo.Create(user)
}

func (uc *useCase) GetUserByID(id string) (*domain.User, error) {
	return uc.userRepo.FindByID(id)
}

func (uc *useCase) GetUserByEmail(email string) (*domain.User, error) {
	return uc.userRepo.FindByEmail(email)
}

func (uc *useCase) UpdateUser(user *domain.User) error {
	return uc.userRepo.Update(user)
}

func (uc *useCase) DeleteUser(id string) error {
	return uc.userRepo.Delete(id)
}

func (uc *useCase) ListUsers() ([]*domain.User, error) {
	return uc.userRepo.List()
}
