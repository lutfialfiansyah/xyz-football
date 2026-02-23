package usecase

import (
	"context"

	"xyz-football-api/internal/modules/match"
	"xyz-football-api/internal/modules/match/repository/postgres"
	playerRepo "xyz-football-api/internal/modules/player/repository/postgres"
	teamRepo "xyz-football-api/internal/modules/team/repository/postgres"
)

type MatchUsecase interface {
	CreateMatch(ctx context.Context, req match.CreateMatchRequest) (*match.Match, error)
	GetAllMatches(ctx context.Context, cursor string, limit int, status string, q string) ([]match.Match, string, bool, error)
	GetMatchByID(ctx context.Context, id string) (*match.Match, error)
	ChangeMatchStatus(ctx context.Context, id string, req match.ChangeMatchStatusRequest) (*match.Match, error)
	ReportMatchScore(ctx context.Context, id string, req match.ReportMatchScoreRequest) (*match.Match, error)
	DeleteMatch(ctx context.Context, id string) error

	AddMatchEvent(ctx context.Context, matchID string, req match.AddMatchEventRequest) (*match.MatchEvent, error)
	GetMatchEvents(ctx context.Context, matchID string) ([]match.MatchEvent, error)
}

type matchUsecase struct {
	matchRepo  postgres.MatchRepository
	teamRepo   teamRepo.TeamRepository
	playerRepo playerRepo.PlayerRepository
}

func NewMatchUsecase(mr postgres.MatchRepository, tr teamRepo.TeamRepository, pr playerRepo.PlayerRepository) MatchUsecase {
	return &matchUsecase{
		matchRepo:  mr,
		teamRepo:   tr,
		playerRepo: pr,
	}
}
