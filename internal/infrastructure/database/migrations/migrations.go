package migrations

import (
	"log"

	"github.com/orinicee/finanzas/internal/domain"
	"gorm.io/gorm"
)

// RunMigrations ejecuta todas las migraciones necesarias
func RunMigrations(db *gorm.DB) error {
	log.Println("Iniciando migraciones...")

	// Auto-migrar todas las entidades
	if err := db.AutoMigrate(
		&domain.User{},
		// Aquí se agregarán más entidades cuando se creen
	); err != nil {
		return err
	}

	log.Println("Migraciones completadas exitosamente")
	return nil
}
