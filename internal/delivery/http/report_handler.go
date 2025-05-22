package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orinicee/finanzas/internal/domain"
)

type ReportHandler struct {
	reportUseCase domain.ReportUseCase
}

func NewReportHandler(reportUseCase domain.ReportUseCase) *ReportHandler {
	return &ReportHandler{
		reportUseCase: reportUseCase,
	}
}

// RegisterRoutes registra las rutas de reportes
func (h *ReportHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	reports := router.Group("/reports")
	reports.Use(authMiddleware)
	{
		reports.GET("/monthly-summary", h.GetMonthlySummary)
		reports.GET("/spending-by-category", h.GetSpendingByCategory)
		reports.GET("/balance-summary", h.GetBalanceSummary)
		reports.GET("/tax-report", h.GetTaxReport)
	}
}

func (h *ReportHandler) GetMonthlySummary(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "año inválido"})
		return
	}

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mes inválido"})
		return
	}

	params := domain.MonthlySummaryParams{
		UserID: userID,
		Year:   year,
		Month:  month,
	}

	summary, err := h.reportUseCase.GenerateMonthlySummary(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *ReportHandler) GetSpendingByCategory(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)

	startDate, err := time.Parse(time.RFC3339, c.Query("start_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha inicial inválida"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha final inválida"})
		return
	}

	params := domain.SpendingByCategoryParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	summary, err := h.reportUseCase.GetSpendingByCategory(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *ReportHandler) GetBalanceSummary(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)

	startDate, err := time.Parse(time.RFC3339, c.Query("start_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha inicial inválida"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha final inválida"})
		return
	}

	groupBy := c.DefaultQuery("group_by", "month")
	if groupBy != "month" && groupBy != "week" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_by debe ser 'month' o 'week'"})
		return
	}

	params := domain.BalanceSummaryParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		GroupBy:   groupBy,
	}

	summary, err := h.reportUseCase.GetBalanceSummary(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *ReportHandler) GetTaxReport(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	userID, _ := uuid.Parse(user.ID)

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "año inválido"})
		return
	}

	params := domain.TaxReportParams{
		UserID: userID,
		Year:   year,
	}

	report, err := h.reportUseCase.GenerateTaxReportForDIAN(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
