package usecase

import (
	"context"
	"net/http"

	"xyz-football-api/internal/modules/match"
	"xyz-football-api/internal/pkg/apperror"
)

func (u *matchUsecase) DeleteMatch(ctx context.Context, id string) error {
	_, err := u.matchRepo.Get(ctx, &match.GetMatchRequest{ID: id})
	if err != nil {
		return apperror.New(http.StatusNotFound, "match not found", "pertandingan tidak ditemukan")
	}
	return u.matchRepo.DeleteMatch(ctx, id)
}
