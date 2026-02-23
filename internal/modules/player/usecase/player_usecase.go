package usecase

import (
	"context"

	"xyz-football-api/internal/modules/player"
	"xyz-football-api/internal/modules/player/repository/postgres"
	teamRepo "xyz-football-api/internal/modules/team/repository/postgres"
)

type PlayerUsecase interface {
	CreatePlayer(ctx context.Context, req player.CreatePlayerRequest) (*player.Player, error)
	GetAllPlayers(ctx context.Context, cursor string, limit int, q string) ([]player.Player, string, bool, error)
	GetPlayerByID(ctx context.Context, id string) (*player.Player, error)
	UpdatePlayer(ctx context.Context, id string, req player.UpdatePlayerRequest) (*player.Player, error)
	DeletePlayer(ctx context.Context, id string) error
}

type playerUsecase struct {
	playerRepo postgres.PlayerRepository
	teamRepo   teamRepo.TeamRepository
}

func NewPlayerUsecase(pr postgres.PlayerRepository, tr teamRepo.TeamRepository) PlayerUsecase {
	return &playerUsecase{
		playerRepo: pr,
		teamRepo:   tr,
	}
}
