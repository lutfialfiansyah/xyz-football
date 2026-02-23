package usecase

import (
	"context"
	"net/http"

	"xyz-football-api/internal/modules/player"
	"xyz-football-api/internal/modules/team"
	"xyz-football-api/internal/pkg/apperror"

	"github.com/google/uuid"
)

func (u *playerUsecase) CreatePlayer(ctx context.Context, req player.CreatePlayerRequest) (*player.Player, error) {
	// Validate Team exists
	_, err := u.teamRepo.Get(ctx, &team.GetTeamRequest{ID: req.TeamID})
	if err != nil {
		return nil, apperror.New(http.StatusBadRequest, "team not found", "tim tidak ditemukan")
	}

	// Validate Jersey Number uniqueness in parsing Team
	existingPlayer, _ := u.playerRepo.Get(ctx, &player.GetPlayerRequest{TeamID: req.TeamID, JerseyNumber: req.JerseyNumber})
	if existingPlayer != nil {
		return nil, apperror.New(http.StatusBadRequest, "jersey number already taken in this team", "nomor jersey sudah dipakai di tim ini")
	}

	teamUUID, _ := uuid.Parse(req.TeamID)

	newPlayer := &player.Player{
		TeamID:       teamUUID,
		Name:         req.Name,
		HeightCm:     req.HeightCm,
		WeightKg:     req.WeightKg,
		Position:     req.Position,
		JerseyNumber: req.JerseyNumber,
	}

	err = u.playerRepo.Create(ctx, newPlayer)
	if err != nil {
		return nil, err
	}

	return newPlayer, nil
}
