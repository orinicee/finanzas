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

	// Verificar si existe la columna password y eliminarla si existe
	var columnExists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'password')").Scan(&columnExists)

	if columnExists {
		if err := db.Exec("ALTER TABLE users DROP COLUMN password").Error; err != nil {
			return err
		}
	}

	log.Println("Migraciones completadas exitosamente")
	return nil
}
