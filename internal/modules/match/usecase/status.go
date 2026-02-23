package usecase

import (
	"context"
	"net/http"

	"xyz-football-api/internal/modules/match"
	"xyz-football-api/internal/pkg/apperror"
)

func (u *matchUsecase) ChangeMatchStatus(ctx context.Context, id string, req match.ChangeMatchStatusRequest) (*match.Match, error) {
	m, err := u.matchRepo.Get(ctx, &match.GetMatchRequest{ID: id})
	if err != nil {
		return nil, apperror.New(http.StatusNotFound, "match not found", "pertandingan tidak ditemukan")
	}

	m.Status = req.Status
	err = u.matchRepo.UpdateMatch(ctx, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (u *matchUsecase) ReportMatchScore(ctx context.Context, id string, req match.ReportMatchScoreRequest) (*match.Match, error) {
	m, err := u.matchRepo.Get(ctx, &match.GetMatchRequest{ID: id})
	if err != nil {
		return nil, apperror.New(http.StatusNotFound, "match not found", "pertandingan tidak ditemukan")
	}

	if m.Status == match.MatchStatusScheduled || m.Status == match.MatchStatusCancelled {
		return nil, apperror.New(http.StatusBadRequest, "cannot report score for scheduled or cancelled match", "tidak bisa melaporkan skor untuk pertandingan yang dijadwalkan atau dibatalkan")
	}

	m.HomeScore = req.HomeScore
	m.AwayScore = req.AwayScore

	err = u.matchRepo.UpdateMatch(ctx, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
