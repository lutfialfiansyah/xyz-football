package postgres

import (
	"context"

	"xyz-football-api/internal/modules/match"
	"xyz-football-api/internal/modules/report"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetFinishedMatches(ctx context.Context) ([]report.MatchRow, error)
	GetGoalEvents(ctx context.Context) ([]report.GoalEvent, error)
}

type reportRepo struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepo{db: db}
}

// GetFinishedMatches returns all finished matches with team names, ordered by date desc.
func (r *reportRepo) GetFinishedMatches(ctx context.Context) ([]report.MatchRow, error) {
	var rows []report.MatchRow

	err := r.db.WithContext(ctx).
		Table("matches m").
		Select(`
			m.id,
			m.match_date,
			m.match_time,
			m.home_team_id,
			m.away_team_id,
			m.home_score,
			m.away_score,
			ht.name as home_team_name,
			at.name as away_team_name
		`).
		Joins("JOIN teams ht ON m.home_team_id = ht.id").
		Joins("JOIN teams at ON m.away_team_id = at.id").
		Where("m.status = ? AND m.deleted_at IS NULL AND ht.deleted_at IS NULL AND at.deleted_at IS NULL", match.MatchStatusFinished).
		Order("m.match_date DESC, m.match_time DESC").
		Scan(&rows).Error

	return rows, err
}

// GetGoalEvents returns all goal events with player names for finished matches.
func (r *reportRepo) GetGoalEvents(ctx context.Context) ([]report.GoalEvent, error) {
	var events []report.GoalEvent

	err := r.db.WithContext(ctx).
		Table("match_events me").
		Select("me.match_id, me.player_id, p.name as player_name").
		Joins("JOIN players p ON me.player_id = p.id").
		Joins("JOIN matches m ON me.match_id = m.id").
		Where("me.event_type = ? AND me.deleted_at IS NULL AND m.status = ? AND m.deleted_at IS NULL", match.EventTypeGoal, match.MatchStatusFinished).
		Scan(&events).Error

	return events, err
}
