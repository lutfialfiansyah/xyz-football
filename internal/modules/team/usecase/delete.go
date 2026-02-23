package usecase

import (
	"context"
	"xyz-football-api/internal/modules/team"
)

func (u *teamUsecase) DeleteTeam(ctx context.Context, id string) error {
	_, err := u.teamRepo.Get(ctx, &team.GetTeamRequest{ID: id})
	if err != nil {
		return err // Team not found
	}

	return u.teamRepo.Delete(ctx, id)
}
