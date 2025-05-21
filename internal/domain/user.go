package domain

import (
	"time"
)

// User representa la entidad de usuario en el dominio
type User struct {
	ID               string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FullName         string     `json:"full_name" gorm:"not null"`
	Email            string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash     string     `json:"-" gorm:"not null"`
	DocumentType     string     `json:"document_type" gorm:"not null"`
	DocumentNumber   string     `json:"document_number" gorm:"not null"`
	TaxRegime        string     `json:"tax_regime" gorm:"not null"`
	PersonType       string     `json:"person_type" gorm:"not null"`
	CIIUCode         string     `json:"ciiu_code"`
	City             string     `json:"city" gorm:"not null"`
	Department       string     `json:"department" gorm:"not null"`
	Address          string     `json:"address" gorm:"not null"`
	Phone            string     `json:"phone" gorm:"not null"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	Subscribed       bool       `json:"subscribed" gorm:"default:false"`
	Role             string     `json:"role" gorm:"not null;default:'user'"`
	HasDianConsent   bool       `json:"has_dian_consent" gorm:"default:false"`
	DianSignatureKey *string    `json:"dian_signature_key"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
