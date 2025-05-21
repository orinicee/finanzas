package domain

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType representa el tipo de transacción
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "ingreso"
	TransactionTypeExpense TransactionType = "egreso"
)

// TransactionCategory representa las categorías de transacciones
type TransactionCategory string

const (
	// Categorías de ingresos
	CategorySalary       TransactionCategory = "salario"
	CategoryFees         TransactionCategory = "honorarios"
	CategoryInterest     TransactionCategory = "intereses"
	CategoryDividends    TransactionCategory = "dividendos"
	CategoryRentalIncome TransactionCategory = "arriendo_recibido"
	CategoryOtherIncome  TransactionCategory = "otros_ingresos"

	// Categorías de egresos
	CategoryRent         TransactionCategory = "arriendo"
	CategoryHealth       TransactionCategory = "salud"
	CategoryEducation    TransactionCategory = "educacion"
	CategoryFood         TransactionCategory = "alimentacion"
	CategoryTransport    TransactionCategory = "transporte"
	CategoryUtilities    TransactionCategory = "servicios_publicos"
	CategoryRecreation   TransactionCategory = "recreacion"
	CategoryOtherExpense TransactionCategory = "otros_egresos"
)

// PaymentMethod representa los métodos de pago
type PaymentMethod string

const (
	PaymentMethodCash      PaymentMethod = "efectivo"
	PaymentMethodDebit     PaymentMethod = "tarjeta_debito"
	PaymentMethodCredit    PaymentMethod = "tarjeta_credito"
	PaymentMethodTransfer  PaymentMethod = "transferencia"
	PaymentMethodNequi     PaymentMethod = "nequi"
	PaymentMethodDaviplata PaymentMethod = "daviplata"
	PaymentMethodOther     PaymentMethod = "otro"
)

// Transaction representa una transacción financiera
type Transaction struct {
	ID              uuid.UUID           `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID           `gorm:"type:uuid;not null" json:"user_id"`
	Type            TransactionType     `gorm:"type:varchar(20);not null" json:"type"`
	Category        TransactionCategory `gorm:"type:varchar(50);not null" json:"category"`
	Amount          float64             `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description     string              `gorm:"type:varchar(255)" json:"description,omitempty"`
	PaymentMethod   PaymentMethod       `gorm:"type:varchar(50);not null" json:"payment_method"`
	Date            time.Time           `gorm:"type:timestamp;not null" json:"date"`
	IsRecurring     bool                `gorm:"type:boolean;default:false" json:"is_recurring"`
	IsTaxDeductible bool                `gorm:"type:boolean;default:false" json:"is_tax_deductible"`
	CreatedAt       time.Time           `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time           `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName especifica el nombre de la tabla para la entidad Transaction
func (Transaction) TableName() string {
	return "transactions"
}
