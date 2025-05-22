package repository

import (
	"context"

	"github.com/orinicee/finanzas/internal/domain"
)

type ReportRepository interface {
	GenerateMonthlySummary(ctx context.Context, params domain.MonthlySummaryParams) (domain.MonthlySummaryResponse, error)
	GetSpendingByCategory(ctx context.Context, params domain.SpendingByCategoryParams) (domain.SpendingByCategoryResponse, error)
	GetBalanceSummary(ctx context.Context, params domain.BalanceSummaryParams) (domain.BalanceSummaryResponse, error)
	GenerateTaxReportForDIAN(ctx context.Context, params domain.TaxReportParams) (domain.TaxReportResponse, error)
}
