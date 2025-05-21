package auth

import (
	"time"

	"github.com/orinicee/finanzas/internal/domain"
	"gorm.io/gorm"
)

// MockUserRepository es un mock del repositorio de usuarios para pruebas
type MockUserRepository struct {
	users         map[string]*domain.User
	refreshTokens map[string]string
}

// NewMockUserRepository crea una nueva instancia del mock
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:         make(map[string]*domain.User),
		refreshTokens: make(map[string]string),
	}
}

func (m *MockUserRepository) Create(user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MockUserRepository) FindBySocialID(provider domain.AuthProvider, socialID string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Provider == provider && user.SocialID == socialID {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MockUserRepository) Update(user *domain.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(id string) error {
	if _, exists := m.users[id]; !exists {
		return ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List() ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) SaveRefreshToken(userID string, token string, expiresAt time.Time) error {
	m.refreshTokens[token] = userID
	return nil
}

func (m *MockUserRepository) GetRefreshToken(token string) (string, error) {
	if userID, exists := m.refreshTokens[token]; exists {
		return userID, nil
	}
	return "", ErrInvalidToken
}

func (m *MockUserRepository) DeleteRefreshToken(token string) error {
	if _, exists := m.refreshTokens[token]; !exists {
		return ErrInvalidToken
	}
	delete(m.refreshTokens, token)
	return nil
}

func (m *MockUserRepository) GetDB() *gorm.DB {
	return nil
}

func (m *MockUserRepository) Close() error {
	return nil
}
