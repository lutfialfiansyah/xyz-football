package usecase

import (
	"context"
	"net/http"
	"time"

	"xyz-football-api/internal/modules/match"
	"xyz-football-api/internal/modules/team"
	"xyz-football-api/internal/pkg/apperror"

	"github.com/google/uuid"
)

func (u *matchUsecase) CreateMatch(ctx context.Context, req match.CreateMatchRequest) (*match.Match, error) {
	if req.HomeTeamID == req.AwayTeamID {
		return nil, apperror.New(http.StatusBadRequest, "home team and away team cannot be the same", "tim tuan rumah dan tim tamu tidak boleh sama")
	}

	_, err := u.teamRepo.Get(ctx, &team.GetTeamRequest{ID: req.HomeTeamID})
	if err != nil {
		return nil, apperror.New(http.StatusBadRequest, "home team not found", "tim tuan rumah tidak ditemukan")
	}

	_, err = u.teamRepo.Get(ctx, &team.GetTeamRequest{ID: req.AwayTeamID})
	if err != nil {
		return nil, apperror.New(http.StatusBadRequest, "away team not found", "tim tamu tidak ditemukan")
	}

	matchDate, err := time.Parse("2006-01-02", req.MatchDate)
	if err != nil {
		return nil, apperror.New(http.StatusBadRequest, "invalid match_date format, expected YYYY-MM-DD", "format tanggal pertandingan tidak valid, gunakan YYYY-MM-DD")
	}

	// Validate match_time format
	_, err = time.Parse("15:04:05", req.MatchTime)
	if err != nil {
		_, errHM := time.Parse("15:04", req.MatchTime)
		if errHM != nil {
			return nil, apperror.New(http.StatusBadRequest, "invalid match_time format, expected HH:MM:SS or HH:MM", "format waktu pertandingan tidak valid, gunakan HH:MM:SS atau HH:MM")
		}
	}

	exists, err := u.matchRepo.CheckMatchScheduleExists(ctx, req.HomeTeamID, req.AwayTeamID, matchDate)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.New(http.StatusBadRequest, "match schedule already exists for these teams on this date", "jadwal pertandingan sudah ada untuk tim-tim ini pada tanggal tersebut")
	}

	homeTeamUUID, _ := uuid.Parse(req.HomeTeamID)
	awayTeamUUID, _ := uuid.Parse(req.AwayTeamID)

	newMatch := &match.Match{
		HomeTeamID: homeTeamUUID,
		AwayTeamID: awayTeamUUID,
		MatchDate:  matchDate,
		MatchTime:  req.MatchTime,
		Status:     match.MatchStatusScheduled,
	}

	err = u.matchRepo.CreateMatch(ctx, newMatch)
	if err != nil {
		return nil, err
	}

	return newMatch, nil
}
