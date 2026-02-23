package usecase

import (
	"context"
	"net/http"

	"xyz-football-api/internal/modules/match"
	"xyz-football-api/internal/modules/player"
	"xyz-football-api/internal/pkg/apperror"

	"github.com/google/uuid"
)

// Events
func (u *matchUsecase) AddMatchEvent(ctx context.Context, matchID string, req match.AddMatchEventRequest) (*match.MatchEvent, error) {
	m, err := u.matchRepo.Get(ctx, &match.GetMatchRequest{ID: matchID})
	if err != nil {
		return nil, apperror.New(http.StatusNotFound, "match not found", "pertandingan tidak ditemukan")
	}

	if m.Status != match.MatchStatusOngoing && m.Status != match.MatchStatusFinished {
		return nil, apperror.New(http.StatusBadRequest, "cannot add event to match that is not ongoing or finished", "tidak bisa menambahkan event ke pertandingan yang belum berlangsung atau selesai")
	}

	p, err := u.playerRepo.Get(ctx, &player.GetPlayerRequest{ID: req.PlayerID})
	if err != nil {
		return nil, apperror.New(http.StatusNotFound, "player not found", "pemain tidak ditemukan")
	}

	// Validate Player belongs to one of the teams in this match
	if p.TeamID != m.HomeTeamID && p.TeamID != m.AwayTeamID {
		return nil, apperror.New(http.StatusBadRequest, "player does not belong to any team playing in this match", "pemain tidak tergabung dalam tim yang bertanding")
	}

	mUUID, _ := uuid.Parse(matchID)
	pUUID, _ := uuid.Parse(req.PlayerID)

	newEvent := &match.MatchEvent{
		MatchID:     mUUID,
		PlayerID:    pUUID,
		TeamID:      p.TeamID, // Auto-assigned from player's current team
		EventMinute: req.EventMinute,
		EventType:   req.EventType,
	}

	err = u.matchRepo.CreateMatchEvent(ctx, newEvent)
	if err != nil {
		return nil, err
	}

	// If event is a GOAL, update the score automatically
	if req.EventType == match.EventTypeGoal {
		if p.TeamID == m.HomeTeamID {
			m.HomeScore += 1
		} else {
			m.AwayScore += 1
		}
		_ = u.matchRepo.UpdateMatch(ctx, m) // Ignore error, best effort
	}
	// If OWN_GOAL, update the opponent's score.
	if req.EventType == match.EventTypeOwnGoal {
		if p.TeamID == m.HomeTeamID {
			m.AwayScore += 1 // Own goal by home team gives point to away team
		} else {
			m.HomeScore += 1
		}
		_ = u.matchRepo.UpdateMatch(ctx, m) // Ignore error, best effort
	}

	return newEvent, nil
}

func (u *matchUsecase) GetMatchEvents(ctx context.Context, matchID string) ([]match.MatchEvent, error) {
	return u.matchRepo.FindEventsByMatchID(ctx, matchID)
}
