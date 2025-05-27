package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/orinicee/finanzas/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrInvalidToken       = errors.New("token inválido")
	ErrTokenExpired       = errors.New("token expirado")
)

type authUseCase struct {
	userRepo domain.UserRepository
	jwtKey   []byte
}

// NewAuthUseCase crea una nueva instancia del caso de uso de autenticación
func NewAuthUseCase(userRepo domain.UserRepository, jwtKey string) domain.AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
	}
}

func (uc *authUseCase) Register(credentials *domain.AuthCredentials) (*domain.AuthToken, error) {
	// Verificar si el usuario ya existe
	existingUser, _ := uc.userRepo.FindByEmail(credentials.Email)
	if existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Crear hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Crear nuevo usuario
	user := &domain.User{
		ID:       uuid.New().String(),
		Email:    credentials.Email,
		Password: string(hashedPassword),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generar tokens
	return uc.generateTokens(user)
}

func (uc *authUseCase) Login(credentials *domain.AuthCredentials) (*domain.AuthToken, error) {
	// Buscar usuario por email
	user, err := uc.userRepo.FindByEmail(credentials.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generar tokens
	return uc.generateTokens(user)
}

// SocialAuth maneja la autenticación con proveedores sociales
// func (uc *authUseCase) SocialAuth(credentials *domain.SocialAuthCredentials) (*domain.AuthToken, error) {
// 	var user *domain.User
// 	var err error

// 	// Validar token según el proveedor
// 	switch credentials.Provider {
// 	case domain.ProviderGoogle:
// 		user, err = uc.validateGoogleToken(credentials.IDToken)
// 	case domain.ProviderApple:
// 		user, err = uc.validateAppleToken(credentials.IDToken)
// 	default:
// 		return nil, errors.New("proveedor no soportado")
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	// Buscar usuario existente o crear uno nuevo
// 	existingUser, err := uc.userRepo.FindBySocialID(credentials.Provider, user.SocialID)
// 	if err != nil {
// 		if err == ErrUserNotFound {
// 			// Crear nuevo usuario
// 			if err := uc.userRepo.Create(user); err != nil {
// 				return nil, err
// 			}
// 			existingUser = user
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	// Generar tokens
// 	return uc.generateTokens(existingUser)
// }

func (uc *authUseCase) RefreshToken(refreshToken string) (*domain.AuthToken, error) {
	// Obtener userID del refresh token
	userID, err := uc.userRepo.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Buscar usuario
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generar nuevos tokens
	return uc.generateTokens(user)
}

func (uc *authUseCase) Logout(refreshToken string) error {
	return uc.userRepo.DeleteRefreshToken(refreshToken)
}

func (uc *authUseCase) ValidateToken(token string) (*domain.User, error) {
	claims := &jwt.RegisteredClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return uc.jwtKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !tkn.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	user, err := uc.userRepo.FindByID(claims.Subject)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// Métodos privados auxiliares

func (uc *authUseCase) generateTokens(user *domain.User) (*domain.AuthToken, error) {
	// Generar access token
	accessTokenExp := time.Now().Add(15 * time.Minute)
	accessTokenClaims := jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(accessTokenExp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(uc.jwtKey)
	if err != nil {
		return nil, err
	}

	// Generar refresh token
	refreshTokenExp := time.Now().Add(7 * 24 * time.Hour) // 7 días
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	refreshTokenString, err := refreshToken.SignedString(uc.jwtKey)
	if err != nil {
		return nil, err
	}

	// Guardar refresh token
	if err := uc.userRepo.SaveRefreshToken(user.ID, refreshTokenString, refreshTokenExp); err != nil {
		return nil, err
	}

	// Calcular expires_in en segundos
	expiresIn := int64(accessTokenExp.Sub(time.Now()).Seconds())

	return &domain.AuthToken{
		AccessToken:  accessTokenString,
		TokenType:    "Bearer",
		ExpiresAt:    accessTokenExp,
		ExpiresIn:    expiresIn,
		RefreshToken: refreshTokenString,
	}, nil
}

// validateGoogleToken valida el token de Google
// func (uc *authUseCase) validateGoogleToken(idToken string) (*domain.User, error) {
// 	// TODO: Implementar validación de token de Google
// 	// Por ahora retornamos un error
// 	return nil, errors.New("validación de Google no implementada")
// }

// validateAppleToken valida el token de Apple
// func (uc *authUseCase) validateAppleToken(idToken string) (*domain.User, error) {
// 	// TODO: Implementar validación de token de Apple
// 	// Por ahora retornamos un error
// 	return nil, errors.New("validación de Apple no implementada")
// }
