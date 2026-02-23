package http

import (
	"net/http"
	"strconv"

	"xyz-football-api/internal/modules/report"
	"xyz-football-api/internal/modules/report/usecase"
	"xyz-football-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportUsecase usecase.ReportUsecase
}

func NewReportHandler(u usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{
		reportUsecase: u,
	}
}

func (h *ReportHandler) GetMatchReports(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "15")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 15
	}

	res, totalData, err := h.reportUsecase.GetMatchReports(c.Request.Context(), page, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	meta := report.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: totalData,
		HasMore:   (page * limit) < totalData,
	}

	response.SuccessWithMeta(c, http.StatusOK, "Match reports fetched successfully", res, meta)
}
