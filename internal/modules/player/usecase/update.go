package usecase

import (
	"context"
	"net/http"

	"xyz-football-api/internal/modules/player"
	"xyz-football-api/internal/modules/team"
	"xyz-football-api/internal/pkg/apperror"

	"github.com/google/uuid"
)

func (u *playerUsecase) UpdatePlayer(ctx context.Context, id string, req player.UpdatePlayerRequest) (*player.Player, error) {
	existingPlayer, err := u.playerRepo.Get(ctx, &player.GetPlayerRequest{ID: id})
	if err != nil {
		return nil, apperror.New(http.StatusBadRequest, "player not found", "pemain tidak ditemukan")
	}

	// Validate Team exists
	if req.TeamID != existingPlayer.TeamID.String() {
		_, err := u.teamRepo.Get(ctx, &team.GetTeamRequest{ID: req.TeamID})
		if err != nil {
			return nil, apperror.New(http.StatusBadRequest, "team not found", "tim tidak ditemukan")
		}
	}

	// Validate Jersey Number if changed
	if req.TeamID != existingPlayer.TeamID.String() || req.JerseyNumber != existingPlayer.JerseyNumber {
		anotherPlayer, _ := u.playerRepo.Get(ctx, &player.GetPlayerRequest{TeamID: req.TeamID, JerseyNumber: req.JerseyNumber})
		if anotherPlayer != nil && anotherPlayer.ID.String() != id {
			return nil, apperror.New(http.StatusBadRequest, "jersey number already taken in this team", "nomor jersey sudah dipakai di tim ini")
		}
	}

	teamUUID, _ := uuid.Parse(req.TeamID)

	existingPlayer.TeamID = teamUUID
	existingPlayer.Name = req.Name
	existingPlayer.HeightCm = req.HeightCm
	existingPlayer.WeightKg = req.WeightKg
	existingPlayer.Position = req.Position
	existingPlayer.JerseyNumber = req.JerseyNumber

	err = u.playerRepo.Update(ctx, existingPlayer)
	if err != nil {
		return nil, err
	}

	return existingPlayer, nil
}
