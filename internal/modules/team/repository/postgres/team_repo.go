package postgres

import (
	"context"
	"time"
	"xyz-football-api/internal/modules/team"

	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(ctx context.Context, t *team.Team) error
	FindAllWithCursor(ctx context.Context, cursorTime time.Time, cursorID string, limit int, q string) ([]team.Team, error)
	Get(ctx context.Context, req *team.GetTeamRequest) (*team.Team, error)
	Update(ctx context.Context, t *team.Team) error
	Delete(ctx context.Context, id string) error
}

type teamRepo struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepo{db: db}
}

func (r *teamRepo) Create(ctx context.Context, t *team.Team) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *teamRepo) FindAllWithCursor(ctx context.Context, cursorTime time.Time, cursorID string, limit int, q string) ([]team.Team, error) {
	var teams []team.Team
	query := r.db.WithContext(ctx).Model(&team.Team{})

	if q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}

	if !cursorTime.IsZero() && cursorID != "" {
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", cursorTime, cursorTime, cursorID)
	}

	// Fetch limit + 1 to determine if there's a next page
	err := query.Order("created_at DESC, id DESC").Limit(limit + 1).Find(&teams).Error
	return teams, err
}

func (r *teamRepo) Get(ctx context.Context, req *team.GetTeamRequest) (*team.Team, error) {
	var output team.Team
	query := r.db.WithContext(ctx).Model(&team.Team{})

	if req.ID != "" {
		query = query.Where("id = ?", req.ID)
	}

	if req.Name != "" {
		query = query.Where("name ILIKE ?", req.Name)
	}

	err := query.First(&output).Error
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (r *teamRepo) Update(ctx context.Context, t *team.Team) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *teamRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&team.Team{}).Error
}
