package report

import (
	"time"

	"github.com/google/uuid"
)

const (
	FinalStatusHomeWin = "Home Win"
	FinalStatusAwayWin = "Away Win"
	FinalStatusDraw    = "Draw"
)

// MatchReportResponse is the final API response DTO.
type MatchReportResponse struct {
	MatchID           string `json:"match_id"`
	MatchDate         string `json:"match_date"`
	MatchTime         string `json:"match_time"`
	HomeTeam          string `json:"home_team"`
	AwayTeam          string `json:"away_team"`
	HomeScore         int    `json:"home_score"`
	AwayScore         int    `json:"away_score"`
	FinalStatus       string `json:"final_status"`
	TopScorer         string `json:"top_scorer"`
	HomeTeamTotalWins int    `json:"home_team_total_wins"`
	AwayTeamTotalWins int    `json:"away_team_total_wins"`
}

// PaginationMeta holds pagination information for the report response.
type PaginationMeta struct {
	Page      int  `json:"page"`
	Limit     int  `json:"limit"`
	TotalData int  `json:"total_data"`
	HasMore   bool `json:"has_more"`
}

// MatchRow is a simple row fetched from the DB (used internally by repo/usecase).
type MatchRow struct {
	ID           uuid.UUID `gorm:"column:id"`
	MatchDate    time.Time `gorm:"column:match_date"`
	MatchTime    string    `gorm:"column:match_time"`
	HomeTeamID   uuid.UUID `gorm:"column:home_team_id"`
	AwayTeamID   uuid.UUID `gorm:"column:away_team_id"`
	HomeScore    int       `gorm:"column:home_score"`
	AwayScore    int       `gorm:"column:away_score"`
	HomeTeamName string    `gorm:"column:home_team_name"`
	AwayTeamName string    `gorm:"column:away_team_name"`
}

// GoalEvent is a single goal event row (used internally by repo/usecase).
type GoalEvent struct {
	MatchID    uuid.UUID `gorm:"column:match_id"`
	PlayerID   uuid.UUID `gorm:"column:player_id"`
	PlayerName string    `gorm:"column:player_name"`
}
