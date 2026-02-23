package usecase

import (
	"context"

	"xyz-football-api/internal/modules/report"
	"xyz-football-api/internal/modules/report/repository/postgres"
)

type ReportUsecase interface {
	GetMatchReports(ctx context.Context, page, limit int) ([]report.MatchReportResponse, int, error)
}

type reportUsecase struct {
	reportRepo postgres.ReportRepository
}

func NewReportUsecase(rr postgres.ReportRepository) ReportUsecase {
	return &reportUsecase{
		reportRepo: rr,
	}
}
