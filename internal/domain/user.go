package domain

import (
	"time"

	"gorm.io/gorm"
)

// User representa un usuario en el sistema
type User struct {
	ID             string         `json:"id" gorm:"primaryKey"`
	FullName       string         `json:"full_name"`
	Email          string         `json:"email" gorm:"uniqueIndex"`
	Password       string         `json:"-" gorm:"not null"` // El "-" evita que se serialice en JSON
	DocumentType   string         `json:"document_type"`
	DocumentNumber string         `json:"document_number"`
	TaxRegime      string         `json:"tax_regime"`
	PersonType     string         `json:"person_type"`
	City           string         `json:"city"`
	Department     string         `json:"department"`
	Address        string         `json:"address"`
	Phone          string         `json:"phone"`
	SocialID       string         `json:"social_id,omitempty" gorm:"uniqueIndex"`
	Provider       AuthProvider   `json:"provider,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}
