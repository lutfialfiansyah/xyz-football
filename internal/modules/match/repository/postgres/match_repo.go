package postgres

import (
	"context"
	"time"

	"xyz-football-api/internal/modules/match"

	"gorm.io/gorm"
)

type MatchRepository interface {
	CreateMatch(ctx context.Context, m *match.Match) error
	FindAllWithCursor(ctx context.Context, cursorDate time.Time, cursorID string, limit int, status string, q string) ([]match.Match, error)
	Get(ctx context.Context, req *match.GetMatchRequest) (*match.Match, error)
	UpdateMatch(ctx context.Context, m *match.Match) error
	DeleteMatch(ctx context.Context, id string) error
	CheckMatchScheduleExists(ctx context.Context, homeTeamID, awayTeamID string, matchDate time.Time) (bool, error)

	CreateMatchEvent(ctx context.Context, event *match.MatchEvent) error
	FindEventsByMatchID(ctx context.Context, matchID string) ([]match.MatchEvent, error)
}

type matchRepo struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepo{db: db}
}

func (r *matchRepo) CreateMatch(ctx context.Context, m *match.Match) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *matchRepo) FindAllWithCursor(ctx context.Context, cursorDate time.Time, cursorID string, limit int, status string, q string) ([]match.Match, error) {
	var matches []match.Match
	query := r.db.WithContext(ctx).Model(&match.Match{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if q != "" {
		// Join with teams to search by team name
		query = query.Joins("LEFT JOIN teams home ON home.id = matches.home_team_id").
			Joins("LEFT JOIN teams away ON away.id = matches.away_team_id").
			Where("home.name ILIKE ? OR away.name ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	if !cursorDate.IsZero() && cursorID != "" {
		query = query.Where("(matches.match_date < ?) OR (matches.match_date = ? AND matches.id < ?)", cursorDate, cursorDate, cursorID)
	}

	// Fetch limit + 1 to determine if there's a next page
	err := query.Order("matches.match_date DESC, matches.id DESC").Limit(limit + 1).Find(&matches).Error
	return matches, err
}

func (r *matchRepo) Get(ctx context.Context, req *match.GetMatchRequest) (*match.Match, error) {
	var m match.Match
	query := r.db.WithContext(ctx).Model(&match.Match{})

	if req.ID != "" {
		query = query.Where("id = ?", req.ID)
	}

	if req.HomeTeamID != "" {
		query = query.Where("home_team_id = ?", req.HomeTeamID)
	}

	if req.AwayTeamID != "" {
		query = query.Where("away_team_id = ?", req.AwayTeamID)
	}

	if !req.MatchDate.IsZero() {
		query = query.Where("match_date = ?", req.MatchDate)
	}

	err := query.First(&m).Error
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *matchRepo) CheckMatchScheduleExists(ctx context.Context, homeTeamID, awayTeamID string, matchDate time.Time) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&match.Match{}).
		Where("home_team_id = ? AND away_team_id = ? AND match_date = ?", homeTeamID, awayTeamID, matchDate).
		Count(&count).Error
	return count > 0, err
}

func (r *matchRepo) UpdateMatch(ctx context.Context, m *match.Match) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *matchRepo) DeleteMatch(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&match.Match{}).Error
}

// Match Events
func (r *matchRepo) CreateMatchEvent(ctx context.Context, event *match.MatchEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *matchRepo) FindEventsByMatchID(ctx context.Context, matchID string) ([]match.MatchEvent, error) {
	var events []match.MatchEvent
	err := r.db.WithContext(ctx).Where("match_id = ?", matchID).Find(&events).Error
	return events, err
}
