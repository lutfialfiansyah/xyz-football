package usecase

import (
	"context"

	"time"

	"xyz-football-api/internal/modules/match"
	"xyz-football-api/internal/pkg/pagination"
)

func (u *matchUsecase) GetAllMatches(ctx context.Context, cursor string, limit int, status string, q string) ([]match.Match, string, bool, error) {
	var cursorDate time.Time
	var cursorID string
	var err error

	if cursor != "" {
		cursorDate, cursorID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return nil, "", false, err
		}
	}

	matches, err := u.matchRepo.FindAllWithCursor(ctx, cursorDate, cursorID, limit, status, q)
	if err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasMore := false

	if len(matches) > limit {
		hasMore = true
		matches = matches[:limit] // Strip the extra item used for checking has_more
	}

	if len(matches) > 0 {
		lastItem := matches[len(matches)-1]
		nextCursor = pagination.EncodeCursor(lastItem.MatchDate, lastItem.ID.String())
	}

	return matches, nextCursor, hasMore, nil
}

func (u *matchUsecase) GetMatchByID(ctx context.Context, id string) (*match.Match, error) {
	return u.matchRepo.Get(ctx, &match.GetMatchRequest{ID: id})
}
