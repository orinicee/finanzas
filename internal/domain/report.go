package domain

import (
	"time"

	"github.com/google/uuid"
)

type MonthlySummaryParams struct {
	UserID uuid.UUID
	Month  int // 1 a 12
	Year   int
}

type MonthlySummaryResponse struct {
	Ingresos         float64
	Egresos          float64
	Balance          float64
	PorcentajeAhorro float64
	Transacciones    int
}

type SpendingByCategoryParams struct {
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
}

type CategorySpending struct {
	Categoria string
	Monto     float64
}

type SpendingByCategoryResponse struct {
	TotalGasto   float64
	PorCategoria []CategorySpending
}

type BalanceSummaryParams struct {
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	GroupBy   string
}

type BalanceSummaryResponse struct {
	Summary []BalanceSummaryItem
}

type BalanceSummaryItem struct {
	Period   string // Ej: "2025-01" o "2025-W03"
	Ingresos float64
	Egresos  float64
	Balance  float64
}

type TaxReportParams struct {
	UserID uuid.UUID
	Year   int
}

type DeductionSummary struct {
	Categoria string
	Monto     float64
}

type TaxReportResponse struct {
	TotalIngresos      float64
	TotalDeducciones   float64
	BaseGravable       float64
	IngresosPorMes     map[string]float64
	DeduccionesDetalle []DeductionSummary
}
