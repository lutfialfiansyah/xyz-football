package usecase

import (
	"context"
	"net/http"

	"xyz-football-api/internal/modules/player"
	"xyz-football-api/internal/pkg/apperror"
)

func (u *playerUsecase) DeletePlayer(ctx context.Context, id string) error {
	_, err := u.playerRepo.Get(ctx, &player.GetPlayerRequest{ID: id})
	if err != nil {
		return apperror.New(http.StatusNotFound, "player not found", "pemain tidak ditemukan")
	}
	return u.playerRepo.Delete(ctx, id)
}
