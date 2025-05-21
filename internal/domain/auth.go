package domain

import "time"

// AuthProvider representa los diferentes proveedores de autenticación
type AuthProvider string

const (
	ProviderLocal  AuthProvider = "local"
	ProviderGoogle AuthProvider = "google"
	ProviderApple  AuthProvider = "apple"
)

// AuthToken representa el token de autenticación
type AuthToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

// AuthCredentials representa las credenciales de autenticación
type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SocialAuthCredentials representa las credenciales de autenticación social
type SocialAuthCredentials struct {
	Provider    AuthProvider `json:"provider"`
	Token       string       `json:"token"`
	IDToken     string       `json:"id_token,omitempty"`
	AccessToken string       `json:"access_token,omitempty"`
}
