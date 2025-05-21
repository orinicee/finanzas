package report

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/orinicee/finanzas/internal/domain"
)

var ()

type ReportUseCase struct {
	transactionRepo domain.TransactionRepository
}

func NewReportUseCase(transactionRepo domain.TransactionRepository) domain.ReportUseCase {
	return &ReportUseCase{
		transactionRepo: transactionRepo,
	}
}

// Genera un resumen financiero del mes, con total de ingresos, egresos y balance.
func (uc *ReportUseCase) GenerateMonthlySummary(ctx context.Context, params domain.MonthlySummaryParams) (domain.MonthlySummaryResponse, error) {
	// Calcular fechas del mes
	startDate := time.Date(params.Year, time.Month(params.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// Obtener transacciones
	transactions, err := uc.transactionRepo.GetByDateRange(ctx, params.UserID, startDate, endDate)
	if err != nil {
		return domain.MonthlySummaryResponse{}, err
	}

	// Calcular totales
	var ingresos, egresos float64
	for _, t := range transactions {
		if t.Type == "ingreso" {
			ingresos += t.Amount
		} else if t.Type == "egreso" {
			egresos += t.Amount
		}
	}

	balance := ingresos - egresos
	var porcentajeAhorro float64
	if ingresos > 0 {
		porcentajeAhorro = (balance / ingresos) * 100
	}

	return domain.MonthlySummaryResponse{
		Ingresos:         ingresos,
		Egresos:          egresos,
		Balance:          balance,
		PorcentajeAhorro: porcentajeAhorro,
		Transacciones:    len(transactions),
	}, nil
}

// Muestra cuánto ha gastado por cada categoría en un rango de fechas.
func (uc *ReportUseCase) GetSpendingByCategory(ctx context.Context, params domain.SpendingByCategoryParams) (domain.SpendingByCategoryResponse, error) {
	transactions, err := uc.transactionRepo.GetByDateRange(ctx, params.UserID, params.StartDate, params.EndDate)
	if err != nil {
		return domain.SpendingByCategoryResponse{}, err
	}

	// Mapa para agrupar gastos por categoría
	categoryTotals := make(map[string]float64)
	var totalGasto float64

	for _, tx := range transactions {
		if tx.Type != "egreso" {
			continue // Solo analizamos egresos
		}
		categoryTotals[string(tx.Category)] += tx.Amount
		totalGasto += tx.Amount
	}

	// Construir respuesta
	var porCategoria []domain.CategorySpending
	for cat, monto := range categoryTotals {
		porCategoria = append(porCategoria, domain.CategorySpending{
			Categoria: cat,
			Monto:     monto,
		})
	}

	return domain.SpendingByCategoryResponse{
		TotalGasto:   totalGasto,
		PorCategoria: porCategoria,
	}, nil
}

// Muestra una vista general del balance histórico del usuario.
func (uc *ReportUseCase) GetBalanceSummary(ctx context.Context, params domain.BalanceSummaryParams) (domain.BalanceSummaryResponse, error) {
	transactions, err := uc.transactionRepo.GetByDateRange(ctx, params.UserID, params.StartDate, params.EndDate)
	if err != nil {
		return domain.BalanceSummaryResponse{}, err
	}

	summaryMap := make(map[string]*domain.BalanceSummaryItem)

	for _, tx := range transactions {
		var periodKey string

		switch params.GroupBy {
		case "week":
			year, week := tx.Date.ISOWeek()
			periodKey = fmt.Sprintf("%04d-W%02d", year, week)
		default: // month
			periodKey = tx.Date.Format("2006-01")
		}

		if _, exists := summaryMap[periodKey]; !exists {
			summaryMap[periodKey] = &domain.BalanceSummaryItem{
				Period: periodKey,
			}
		}

		if tx.Type == domain.TransactionTypeIncome {
			summaryMap[periodKey].Ingresos += tx.Amount
		} else if tx.Type == domain.TransactionTypeExpense {
			summaryMap[periodKey].Egresos += tx.Amount
		}
	}

	// Convertir el mapa en un slice
	var result []domain.BalanceSummaryItem
	for _, item := range summaryMap {
		item.Balance = item.Ingresos - item.Egresos
		result = append(result, *item)
	}

	// (Opcional) ordenar los resultados cronológicamente
	sort.Slice(result, func(i, j int) bool {
		return result[i].Period < result[j].Period
	})

	return domain.BalanceSummaryResponse{Summary: result}, nil
}

// Genera un reporte tributario con los datos relevantes para declarar ante la DIAN.
func (uc *ReportUseCase) GenerateTaxReportForDIAN(ctx context.Context, params domain.TaxReportParams) (domain.TaxReportResponse, error) {
	startDate := time.Date(params.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(params.Year, 12, 31, 23, 59, 59, 999999999, time.UTC)

	transactions, err := uc.transactionRepo.GetByDateRange(ctx, params.UserID, startDate, endDate)
	if err != nil {
		return domain.TaxReportResponse{}, err
	}

	var totalIngresos float64
	var totalDeducciones float64
	ingresosPorMes := make(map[string]float64)
	deduccionesPorCategoria := make(map[string]float64)

	for _, tx := range transactions {
		month := tx.Date.Format("01") // "01" = mes en dos dígitos
		if tx.Type == domain.TransactionTypeIncome {
			totalIngresos += tx.Amount
			ingresosPorMes[month] += tx.Amount
		} else if tx.Type == domain.TransactionTypeExpense && tx.IsTaxDeductible {
			totalDeducciones += tx.Amount
			deduccionesPorCategoria[string(tx.Category)] += tx.Amount
		}
	}

	// Convertir deducciones a lista
	var deduccionesDetalle []domain.DeductionSummary
	for categoria, monto := range deduccionesPorCategoria {
		deduccionesDetalle = append(deduccionesDetalle, domain.DeductionSummary{
			Categoria: categoria,
			Monto:     monto,
		})
	}

	resp := domain.TaxReportResponse{
		TotalIngresos:      totalIngresos,
		TotalDeducciones:   totalDeducciones,
		BaseGravable:       totalIngresos - totalDeducciones,
		IngresosPorMes:     ingresosPorMes,
		DeduccionesDetalle: deduccionesDetalle,
	}

	return resp, nil
}
