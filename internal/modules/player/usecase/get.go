package usecase

import (
	"context"

	"time"

	"xyz-football-api/internal/modules/player"
	"xyz-football-api/internal/pkg/pagination"
)

func (u *playerUsecase) GetAllPlayers(ctx context.Context, cursor string, limit int, q string) ([]player.Player, string, bool, error) {
	var cursorTime time.Time
	var cursorID string
	var err error

	if cursor != "" {
		cursorTime, cursorID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return nil, "", false, err
		}
	}

	players, err := u.playerRepo.FindAllWithCursor(ctx, cursorTime, cursorID, limit, q)
	if err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasMore := false

	if len(players) > limit {
		hasMore = true
		players = players[:limit] // Strip the extra item used for checking has_more
	}

	if len(players) > 0 {
		lastItem := players[len(players)-1]
		nextCursor = pagination.EncodeCursor(lastItem.CreatedAt, lastItem.ID.String())
	}

	return players, nextCursor, hasMore, nil
}

func (u *playerUsecase) GetPlayerByID(ctx context.Context, id string) (*player.Player, error) {
	return u.playerRepo.Get(ctx, &player.GetPlayerRequest{ID: id})
}
