package match

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Match Status Enumerations
const (
	MatchStatusScheduled = "scheduled"
	MatchStatusOngoing   = "ongoing"
	MatchStatusFinished  = "finished"
	MatchStatusCancelled = "cancelled"
)

// Match Event Type Enumerations
const (
	EventTypeGoal       = "goal"
	EventTypeYellowCard = "yellow_card"
	EventTypeRedCard    = "red_card"
	EventTypeOwnGoal    = "own_goal"
)

type Match struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	HomeTeamID uuid.UUID      `gorm:"type:uuid;not null;index:idx_matches_home_team,where:deleted_at IS NULL" json:"home_team_id"`
	AwayTeamID uuid.UUID      `gorm:"type:uuid;not null;index:idx_matches_away_team,where:deleted_at IS NULL" json:"away_team_id"`
	MatchDate  time.Time      `gorm:"type:date;not null;index:idx_matches_date,where:deleted_at IS NULL;uniqueIndex:idx_unique_match_schedule,where:deleted_at IS NULL" json:"match_date"`
	MatchTime  string         `gorm:"type:time;not null" json:"match_time"`
	HomeScore  int            `gorm:"type:integer;default:0" json:"home_score"`
	AwayScore  int            `gorm:"type:integer;default:0" json:"away_score"`
	Status     string         `gorm:"type:varchar(20);default:'scheduled';index:idx_matches_status,where:deleted_at IS NULL" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type MatchEvent struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MatchID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_events_match,where:deleted_at IS NULL" json:"match_id"`
	PlayerID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_events_player,where:deleted_at IS NULL" json:"player_id"`
	TeamID      uuid.UUID      `gorm:"type:uuid;not null" json:"team_id"`
	EventMinute int            `gorm:"type:integer;not null" json:"event_minute"`
	EventType   string         `gorm:"type:varchar(20);not null;index:idx_events_type,where:deleted_at IS NULL" json:"event_type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Custom parser to handle date parsing
type CreateMatchRequest struct {
	HomeTeamID string `json:"home_team_id" binding:"required,uuid"`
	AwayTeamID string `json:"away_team_id" binding:"required,uuid"`
	MatchDate  string `json:"match_date" binding:"required"` // Format: YYYY-MM-DD
	MatchTime  string `json:"match_time" binding:"required"` // Format: HH:MM:SS
}

type ChangeMatchStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=scheduled ongoing finished cancelled"`
}

type ReportMatchScoreRequest struct {
	HomeScore int `json:"home_score" binding:"gte=0"`
	AwayScore int `json:"away_score" binding:"gte=0"`
}

type AddMatchEventRequest struct {
	PlayerID    string `json:"player_id" binding:"required,uuid"`
	EventMinute int    `json:"event_minute" binding:"required,gt=0,lte=120"`
	EventType   string `json:"event_type" binding:"required,oneof=goal yellow_card red_card own_goal"`
}

type GetMatchRequest struct {
	ID         string
	HomeTeamID string
	AwayTeamID string
	MatchDate  time.Time
}
