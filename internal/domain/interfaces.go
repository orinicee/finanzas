package domain

import (
	"time"
)

// UserRepository define la interfaz para el repositorio de usuarios
type UserRepository interface {
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
