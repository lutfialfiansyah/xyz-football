package postgres

import (
	"context"
	"time"

	"xyz-football-api/internal/modules/player"

	"gorm.io/gorm"
)

type PlayerRepository interface {
	Create(ctx context.Context, p *player.Player) error
	FindAllWithCursor(ctx context.Context, cursorTime time.Time, cursorID string, limit int, q string) ([]player.Player, error)
	Get(ctx context.Context, req *player.GetPlayerRequest) (*player.Player, error)
	Update(ctx context.Context, p *player.Player) error
	Delete(ctx context.Context, id string) error
}

type playerRepo struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepo{db: db}
}

func (r *playerRepo) Create(ctx context.Context, p *player.Player) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *playerRepo) FindAllWithCursor(ctx context.Context, cursorTime time.Time, cursorID string, limit int, q string) ([]player.Player, error) {
	var players []player.Player
	query := r.db.WithContext(ctx).Model(&player.Player{})

	if q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}

	if !cursorTime.IsZero() && cursorID != "" {
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", cursorTime, cursorTime, cursorID)
	}

	// Fetch limit + 1 to determine if there's a next page
	err := query.Order("created_at DESC, id DESC").Limit(limit + 1).Find(&players).Error
	return players, err
}

func (r *playerRepo) Get(ctx context.Context, req *player.GetPlayerRequest) (*player.Player, error) {
	var output player.Player
	query := r.db.WithContext(ctx).Model(&player.Player{})

	if req.ID != "" {
		query = query.Where("id = ?", req.ID)
	}

	if req.TeamID != "" {
		query = query.Where("team_id = ?", req.TeamID)
	}

	if req.JerseyNumber > 0 {
		query = query.Where("jersey_number = ?", req.JerseyNumber)
	}

	err := query.First(&output).Error
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (r *playerRepo) Update(ctx context.Context, p *player.Player) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *playerRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&player.Player{}).Error
}
