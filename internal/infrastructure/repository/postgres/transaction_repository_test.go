package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/orinicee/finanzas/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Obtener configuración de variables de entorno o usar valores por defecto
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres123")
	dbname := getEnv("DB_NAME", "finanzas_test")

	// Construir DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Habilitar la extensión uuid-ossp
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	require.NoError(t, err)

	// Limpiar la base de datos antes de cada test
	err = db.Migrator().DropTable(&domain.Transaction{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Transaction{})
	require.NoError(t, err)

	return db
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func TestTransactionRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	// Crear una transacción de prueba
	transaction := &domain.Transaction{
		UserID:          uuid.New(),
		Type:            domain.TransactionTypeIncome,
		Category:        domain.CategorySalary,
		Amount:          1000000.00,
		Description:     "Salario mensual",
		PaymentMethod:   domain.PaymentMethodTransfer,
		Date:            time.Now(),
		IsRecurring:     true,
		IsTaxDeductible: false,
	}

	// Probar la creación
	err := repo.Create(ctx, transaction)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, transaction.ID)

	// Verificar que se puede recuperar
	retrieved, err := repo.GetByID(ctx, transaction.ID)
	require.NoError(t, err)
	assert.Equal(t, transaction.UserID, retrieved.UserID)
	assert.Equal(t, transaction.Amount, retrieved.Amount)
}

func TestTransactionRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	userID := uuid.New()

	// Crear múltiples transacciones para el mismo usuario
	transactions := []*domain.Transaction{
		{
			UserID:        userID,
			Type:          domain.TransactionTypeIncome,
			Category:      domain.CategorySalary,
			Amount:        1000000.00,
			PaymentMethod: domain.PaymentMethodTransfer,
			Date:          time.Now(),
		},
		{
			UserID:        userID,
			Type:          domain.TransactionTypeExpense,
			Category:      domain.CategoryRent,
			Amount:        500000.00,
			PaymentMethod: domain.PaymentMethodTransfer,
			Date:          time.Now(),
		},
	}

	// Crear las transacciones
	for _, tx := range transactions {
		err := repo.Create(ctx, tx)
		require.NoError(t, err)
	}

	// Obtener todas las transacciones del usuario
	retrieved, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, retrieved, 2)
}

func TestTransactionRepository_GetByDateRange(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	now := time.Now()

	// Crear transacciones en diferentes fechas
	transactions := []*domain.Transaction{
		{
			UserID:        userID,
			Type:          domain.TransactionTypeIncome,
			Category:      domain.CategorySalary,
			Amount:        1000000.00,
			PaymentMethod: domain.PaymentMethodTransfer,
			Date:          now.AddDate(0, -1, 0), // Hace un mes
		},
		{
			UserID:        userID,
			Type:          domain.TransactionTypeExpense,
			Category:      domain.CategoryRent,
			Amount:        500000.00,
			PaymentMethod: domain.PaymentMethodTransfer,
			Date:          now, // Hoy
		},
	}

	// Crear las transacciones
	for _, tx := range transactions {
		err := repo.Create(ctx, tx)
		require.NoError(t, err)
	}

	// Obtener transacciones del último mes
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 0, 1)
	retrieved, err := repo.GetByDateRange(ctx, userID, startDate, endDate)
	require.NoError(t, err)
	assert.Len(t, retrieved, 2)
}

func TestTransactionRepository_GetByCategory(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	userID := uuid.New()

	// Crear transacciones de diferentes categorías
	transactions := []*domain.Transaction{
		{
			UserID:        userID,
			Type:          domain.TransactionTypeIncome,
			Category:      domain.CategorySalary,
			Amount:        1000000.00,
			PaymentMethod: domain.PaymentMethodTransfer,
			Date:          time.Now(),
		},
		{
			UserID:        userID,
			Type:          domain.TransactionTypeExpense,
			Category:      domain.CategoryRent,
			Amount:        500000.00,
			PaymentMethod: domain.PaymentMethodTransfer,
			Date:          time.Now(),
		},
	}

	// Crear las transacciones
	for _, tx := range transactions {
		err := repo.Create(ctx, tx)
		require.NoError(t, err)
	}

	// Obtener transacciones por categoría
	retrieved, err := repo.GetByCategory(ctx, userID, domain.CategorySalary)
	require.NoError(t, err)
	assert.Len(t, retrieved, 1)
	assert.Equal(t, domain.CategorySalary, retrieved[0].Category)
}

func TestTransactionRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	// Crear una transacción
	transaction := &domain.Transaction{
		UserID:        uuid.New(),
		Type:          domain.TransactionTypeIncome,
		Category:      domain.CategorySalary,
		Amount:        1000000.00,
		PaymentMethod: domain.PaymentMethodTransfer,
		Date:          time.Now(),
	}

	err := repo.Create(ctx, transaction)
	require.NoError(t, err)

	// Modificar la transacción
	transaction.Amount = 1500000.00
	transaction.Description = "Actualización de salario"

	// Actualizar la transacción
	err = repo.Update(ctx, transaction)
	require.NoError(t, err)

	// Verificar la actualización
	retrieved, err := repo.GetByID(ctx, transaction.ID)
	require.NoError(t, err)
	assert.Equal(t, 1500000.00, retrieved.Amount)
	assert.Equal(t, "Actualización de salario", retrieved.Description)
}

func TestTransactionRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	// Crear una transacción
	transaction := &domain.Transaction{
		UserID:        uuid.New(),
		Type:          domain.TransactionTypeIncome,
		Category:      domain.CategorySalary,
		Amount:        1000000.00,
		PaymentMethod: domain.PaymentMethodTransfer,
		Date:          time.Now(),
	}

	err := repo.Create(ctx, transaction)
	require.NoError(t, err)

	// Eliminar la transacción
	err = repo.Delete(ctx, transaction.ID)
	require.NoError(t, err)

	// Verificar que no existe
	_, err = repo.GetByID(ctx, transaction.ID)
	assert.Error(t, err)
}
