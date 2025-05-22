package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orinicee/finanzas/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReportUseCase struct {
	mock.Mock
}

func (m *MockReportUseCase) GenerateMonthlySummary(ctx context.Context, params domain.MonthlySummaryParams) (domain.MonthlySummaryResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(domain.MonthlySummaryResponse), args.Error(1)
}

func (m *MockReportUseCase) GetSpendingByCategory(ctx context.Context, params domain.SpendingByCategoryParams) (domain.SpendingByCategoryResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(domain.SpendingByCategoryResponse), args.Error(1)
}

func (m *MockReportUseCase) GetBalanceSummary(ctx context.Context, params domain.BalanceSummaryParams) (domain.BalanceSummaryResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(domain.BalanceSummaryResponse), args.Error(1)
}

func (m *MockReportUseCase) GenerateTaxReportForDIAN(ctx context.Context, params domain.TaxReportParams) (domain.TaxReportResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(domain.TaxReportResponse), args.Error(1)
}

func setupTestRouters(handler *ReportHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(&router.RouterGroup, func(c *gin.Context) {
		// Mock user for testing
		user := &domain.User{
			ID: uuid.New().String(),
		}
		c.Set("user", user)
		c.Next()
	})
	return router
}

func TestGetMonthlySummary(t *testing.T) {
	mockUseCase := new(MockReportUseCase)
	handler := NewReportHandler(mockUseCase)
	router := setupTestRouters(handler)

	expectedResponse := domain.MonthlySummaryResponse{
		Ingresos:         1000,
		Egresos:          500,
		Balance:          500,
		PorcentajeAhorro: 50,
		Transacciones:    10,
	}

	mockUseCase.On("GenerateMonthlySummary", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/reports/monthly-summary?year=2024&month=3", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.MonthlySummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestGetSpendingByCategory(t *testing.T) {
	mockUseCase := new(MockReportUseCase)
	handler := NewReportHandler(mockUseCase)
	router := setupTestRouters(handler)

	expectedResponse := domain.SpendingByCategoryResponse{
		TotalGasto: 1000,
		PorCategoria: []domain.CategorySpending{
			{
				Categoria: "Alimentación",
				Monto:     500,
			},
			{
				Categoria: "Transporte",
				Monto:     500,
			},
		},
	}

	mockUseCase.On("GetSpendingByCategory", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	startDate := time.Now().AddDate(0, -1, 0).Format(time.RFC3339)
	endDate := time.Now().Format(time.RFC3339)
	req, _ := http.NewRequest("GET", "/api/reports/spending-by-category?start_date="+startDate+"&end_date="+endDate, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.SpendingByCategoryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestGetBalanceSummary(t *testing.T) {
	mockUseCase := new(MockReportUseCase)
	handler := NewReportHandler(mockUseCase)
	router := setupTestRouters(handler)

	expectedResponse := domain.BalanceSummaryResponse{
		Summary: []domain.BalanceSummaryItem{
			{
				Period:   "2024-01",
				Ingresos: 1000,
				Egresos:  500,
				Balance:  500,
			},
		},
	}

	mockUseCase.On("GetBalanceSummary", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	startDate := time.Now().AddDate(0, -1, 0).Format(time.RFC3339)
	endDate := time.Now().Format(time.RFC3339)
	req, _ := http.NewRequest("GET", "/api/reports/balance-summary?start_date="+startDate+"&end_date="+endDate+"&group_by=month", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.BalanceSummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestGetTaxReport(t *testing.T) {
	mockUseCase := new(MockReportUseCase)
	handler := NewReportHandler(mockUseCase)
	router := setupTestRouters(handler)

	expectedResponse := domain.TaxReportResponse{
		TotalIngresos:    10000,
		TotalDeducciones: 2000,
		BaseGravable:     8000,
		IngresosPorMes: map[string]float64{
			"01": 1000,
			"02": 1000,
		},
		DeduccionesDetalle: []domain.DeductionSummary{
			{
				Categoria: "Salud",
				Monto:     1000,
			},
			{
				Categoria: "Educación",
				Monto:     1000,
			},
		},
	}

	mockUseCase.On("GenerateTaxReportForDIAN", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/reports/tax-report?year=2024", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.TaxReportResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}
