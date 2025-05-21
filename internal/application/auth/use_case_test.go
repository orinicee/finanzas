package auth

import (
	"testing"
	"time"

	"github.com/orinicee/finanzas/internal/domain"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name          string
		credentials   *domain.AuthCredentials
		setupMock     func(*MockUserRepository)
		expectedError error
	}{
		{
			name: "registro exitoso",
			credentials: &domain.AuthCredentials{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				// No setup needed for successful registration
			},
			expectedError: nil,
		},
		{
			name: "email ya registrado",
			credentials: &domain.AuthCredentials{
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				m.Create(&domain.User{
					ID:       "1",
					Email:    "existing@example.com",
					Password: "hashed_password",
				})
			},
			expectedError: domain.ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockUserRepository()
			tt.setupMock(mockRepo)
			useCase := NewAuthUseCase(mockRepo, "test-key")

			// Execute
			token, err := useCase.Register(tt.credentials)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.RefreshToken)
				assert.Equal(t, "Bearer", token.TokenType)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name          string
		credentials   *domain.AuthCredentials
		setupMock     func(*MockUserRepository)
		expectedError error
	}{
		{
			name: "login exitoso",
			credentials: &domain.AuthCredentials{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				m.Create(&domain.User{
					ID:       "1",
					Email:    "test@example.com",
					Password: string(hashedPassword),
				})
			},
			expectedError: nil,
		},
		{
			name: "credenciales inválidas",
			credentials: &domain.AuthCredentials{
				Email:    "test@example.com",
				Password: "wrong_password",
			},
			setupMock: func(m *MockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				m.Create(&domain.User{
					ID:       "1",
					Email:    "test@example.com",
					Password: string(hashedPassword),
				})
			},
			expectedError: ErrInvalidCredentials,
		},
		{
			name: "usuario no encontrado",
			credentials: &domain.AuthCredentials{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock:     func(m *MockUserRepository) {},
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockUserRepository()
			tt.setupMock(mockRepo)
			useCase := NewAuthUseCase(mockRepo, "test-key")

			// Execute
			token, err := useCase.Login(tt.credentials)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.RefreshToken)
				assert.Equal(t, "Bearer", token.TokenType)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name          string
		refreshToken  string
		setupMock     func(*MockUserRepository)
		expectedError error
	}{
		{
			name:         "logout exitoso",
			refreshToken: "valid-token",
			setupMock: func(m *MockUserRepository) {
				m.SaveRefreshToken("user1", "valid-token", time.Now().Add(time.Hour))
			},
			expectedError: nil,
		},
		{
			name:          "token inválido",
			refreshToken:  "invalid-token",
			setupMock:     func(m *MockUserRepository) {},
			expectedError: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockUserRepository()
			tt.setupMock(mockRepo)
			useCase := NewAuthUseCase(mockRepo, "test-key")

			// Execute
			err := useCase.Logout(tt.refreshToken)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				// Verificar que el token fue eliminado
				_, err := mockRepo.GetRefreshToken(tt.refreshToken)
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidToken, err)
			}
		})
	}
}
