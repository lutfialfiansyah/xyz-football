package usecase

import (
	"context"

	"xyz-football-api/internal/modules/team"
	"xyz-football-api/internal/modules/team/repository/postgres"
	"xyz-football-api/internal/pkg/storage"
)

type TeamUsecase interface {
	CreateTeam(ctx context.Context, req team.CreateTeamRequest) (*team.Team, error)
	GetAllTeams(ctx context.Context, cursor string, limit int, q string) ([]team.Team, string, bool, error)
	GetTeamByID(ctx context.Context, id string) (*team.Team, error)
	UpdateTeam(ctx context.Context, id string, req team.UpdateTeamRequest) (*team.Team, error)
	DeleteTeam(ctx context.Context, id string) error
}

type teamUsecase struct {
	teamRepo postgres.TeamRepository
	storage  storage.Provider
}

func NewTeamUsecase(repo postgres.TeamRepository, s storage.Provider) TeamUsecase {
	return &teamUsecase{
		teamRepo: repo,
		storage:  s,
	}
}
