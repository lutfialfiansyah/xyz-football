package usecase

import (
	"context"

	"time"

	"xyz-football-api/internal/modules/team"
	"xyz-football-api/internal/pkg/pagination"
)

func (u *teamUsecase) GetAllTeams(ctx context.Context, cursor string, limit int, q string) ([]team.Team, string, bool, error) {
	var cursorTime time.Time
	var cursorID string
	var err error

	if cursor != "" {
		cursorTime, cursorID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return nil, "", false, err
		}
	}

	teams, err := u.teamRepo.FindAllWithCursor(ctx, cursorTime, cursorID, limit, q)
	if err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasMore := false

	if len(teams) > limit {
		hasMore = true
		teams = teams[:limit] // Strip the extra item used for checking has_more
	}

	if len(teams) > 0 {
		lastItem := teams[len(teams)-1]
		nextCursor = pagination.EncodeCursor(lastItem.CreatedAt, lastItem.ID.String())
	}

	// Map LogoURL to absolute URL
	for i := range teams {
		teams[i].LogoURL = u.storage.GetURL(teams[i].LogoURL)
	}

	return teams, nextCursor, hasMore, nil
}

func (u *teamUsecase) GetTeamByID(ctx context.Context, id string) (*team.Team, error) {
	teamData, err := u.teamRepo.Get(ctx, &team.GetTeamRequest{ID: id})
	if err != nil {
		return nil, err
	}

	teamData.LogoURL = u.storage.GetURL(teamData.LogoURL)
	return teamData, nil
}
