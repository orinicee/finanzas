package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/orinicee/finanzas/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserUseCase es un mock del caso de uso de usuario
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserUseCase) GetUserByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUseCase) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserUseCase) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserUseCase) ListUsers() ([]*domain.User, error) {
	args := m.Called()
	return args.Get(0).([]*domain.User), args.Error(1)
}

// MockUserRepository es un mock del repositorio de usuario
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List() ([]*domain.User, error) {
	args := m.Called()
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.DB)
}

func (m *MockUserRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		mockSetup      func(*MockUserUseCase, *MockUserRepository)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "crear usuario exitosamente",
			payload: map[string]interface{}{
				"full_name":       "Test User",
				"email":           "test@example.com",
				"password":        "password123",
				"document_type":   "CC",
				"document_number": "123456789",
				"tax_regime":      "Común",
				"person_type":     "Natural",
				"city":            "Bogotá",
				"department":      "Cundinamarca",
				"address":         "Calle 123",
				"phone":           "1234567890",
			},
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				uc.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Usuario creado exitosamente",
			},
		},
		{
			name: "error al crear usuario - email duplicado",
			payload: map[string]interface{}{
				"full_name":       "Test User",
				"email":           "existing@example.com",
				"password":        "password123",
				"document_type":   "CC",
				"document_number": "123456789",
				"tax_regime":      "Común",
				"person_type":     "Natural",
				"city":            "Bogotá",
				"department":      "Cundinamarca",
				"address":         "Calle 123",
				"phone":           "1234567890",
			},
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				uc.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(domain.ErrEmailAlreadyExists)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "El email ya está registrado",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			mockUseCase := new(MockUserUseCase)
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockUseCase, mockRepo)

			// Crear handler
			handler := NewUserHandler(mockUseCase, mockRepo)

			// Configurar router de prueba
			router := setupTestRouter()
			router.POST("/users", handler.CreateUser)

			// Crear request de prueba
			jsonData, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Ejecutar request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificar resultados
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Verificar que los mocks fueron llamados
			mockUseCase.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserUseCase, *MockUserRepository)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "obtener usuario exitosamente",
			userID: "1",
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				uc.On("GetUserByID", "1").Return(&domain.User{
					ID:             "1",
					FullName:       "Test User",
					Email:          "test@example.com",
					DocumentType:   "CC",
					DocumentNumber: "123456789",
					TaxRegime:      "Común",
					PersonType:     "Natural",
					City:           "Bogotá",
					Department:     "Cundinamarca",
					Address:        "Calle 123",
					Phone:          "1234567890",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":              "1",
				"full_name":       "Test User",
				"email":           "test@example.com",
				"document_type":   "CC",
				"document_number": "123456789",
				"tax_regime":      "Común",
				"person_type":     "Natural",
				"city":            "Bogotá",
				"department":      "Cundinamarca",
				"address":         "Calle 123",
				"phone":           "1234567890",
			},
		},
		{
			name:   "usuario no encontrado",
			userID: "999",
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				uc.On("GetUserByID", "999").Return(nil, domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Usuario no encontrado",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			mockUseCase := new(MockUserUseCase)
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockUseCase, mockRepo)

			// Crear handler
			handler := NewUserHandler(mockUseCase, mockRepo)

			// Configurar router de prueba
			router := setupTestRouter()
			router.GET("/users/:id", handler.GetUser)

			// Crear request de prueba
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)

			// Ejecutar request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificar resultados
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Verificar que los mocks fueron llamados
			mockUseCase.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		payload        map[string]interface{}
		mockSetup      func(*MockUserUseCase, *MockUserRepository)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "actualizar usuario exitosamente",
			userID: "1",
			payload: map[string]interface{}{
				"full_name": "Updated User",
				"email":     "updated@example.com",
			},
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				uc.On("UpdateUser", mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Usuario actualizado exitosamente",
			},
		},
		{
			name:   "error al actualizar - usuario no encontrado",
			userID: "999",
			payload: map[string]interface{}{
				"full_name": "Updated User",
				"email":     "updated@example.com",
			},
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				uc.On("UpdateUser", mock.AnythingOfType("*domain.User")).Return(domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Usuario no encontrado",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			mockUseCase := new(MockUserUseCase)
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockUseCase, mockRepo)

			// Crear handler
			handler := NewUserHandler(mockUseCase, mockRepo)

			// Configurar router de prueba
			router := setupTestRouter()
			router.PUT("/users/:id", handler.UpdateUser)

			// Crear request de prueba
			jsonData, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userID, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Ejecutar request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificar resultados
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Verificar que los mocks fueron llamados
			mockUseCase.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserUseCase, *MockUserRepository)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "eliminar usuario exitosamente",
			userID: "1",
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				uc.On("DeleteUser", "1").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Usuario eliminado exitosamente",
			},
		},
		{
			name:   "error al eliminar - usuario no encontrado",
			userID: "999",
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				uc.On("DeleteUser", "999").Return(domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Usuario no encontrado",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			mockUseCase := new(MockUserUseCase)
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockUseCase, mockRepo)

			// Crear handler
			handler := NewUserHandler(mockUseCase, mockRepo)

			// Configurar router de prueba
			router := setupTestRouter()
			router.DELETE("/users/:id", handler.DeleteUser)

			// Crear request de prueba
			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)

			// Ejecutar request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificar resultados
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Verificar que los mocks fueron llamados
			mockUseCase.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestListUsers(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockUserUseCase, *MockUserRepository)
		expectedStatus int
		expectedBody   []map[string]interface{}
	}{
		{
			name: "listar usuarios exitosamente",
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				users := []*domain.User{
					{
						ID:             "1",
						FullName:       "User 1",
						Email:          "user1@example.com",
						DocumentType:   "CC",
						DocumentNumber: "123456789",
						TaxRegime:      "Común",
						PersonType:     "Natural",
						City:           "Bogotá",
						Department:     "Cundinamarca",
						Address:        "Calle 123",
						Phone:          "1234567890",
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					},
					{
						ID:             "2",
						FullName:       "User 2",
						Email:          "user2@example.com",
						DocumentType:   "CC",
						DocumentNumber: "987654321",
						TaxRegime:      "Común",
						PersonType:     "Natural",
						City:           "Medellín",
						Department:     "Antioquia",
						Address:        "Calle 456",
						Phone:          "0987654321",
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					},
				}
				uc.On("ListUsers").Return(users, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []map[string]interface{}{
				{
					"id":              "1",
					"full_name":       "User 1",
					"email":           "user1@example.com",
					"document_type":   "CC",
					"document_number": "123456789",
					"tax_regime":      "Común",
					"person_type":     "Natural",
					"city":            "Bogotá",
					"department":      "Cundinamarca",
					"address":         "Calle 123",
					"phone":           "1234567890",
				},
				{
					"id":              "2",
					"full_name":       "User 2",
					"email":           "user2@example.com",
					"document_type":   "CC",
					"document_number": "987654321",
					"tax_regime":      "Común",
					"person_type":     "Natural",
					"city":            "Medellín",
					"department":      "Antioquia",
					"address":         "Calle 456",
					"phone":           "0987654321",
				},
			},
		},
		{
			name: "lista vacía de usuarios",
			mockSetup: func(uc *MockUserUseCase, repo *MockUserRepository) {
				uc.On("ListUsers").Return([]*domain.User{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mocks
			mockUseCase := new(MockUserUseCase)
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockUseCase, mockRepo)

			// Crear handler
			handler := NewUserHandler(mockUseCase, mockRepo)

			// Configurar router de prueba
			router := setupTestRouter()
			router.GET("/users", handler.ListUsers)

			// Crear request de prueba
			req := httptest.NewRequest(http.MethodGet, "/users", nil)

			// Ejecutar request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificar resultados
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Verificar que los mocks fueron llamados
			mockUseCase.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
