package player

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Player struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TeamID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_players_team,where:deleted_at IS NULL" json:"team_id"`
	Name         string         `gorm:"type:varchar(100);not null;index:idx_players_name,where:deleted_at IS NULL" json:"name"`
	HeightCm     float64        `gorm:"type:numeric(5,2)" json:"height_cm"`
	WeightKg     float64        `gorm:"type:numeric(5,2)" json:"weight_kg"`
	Position     string         `gorm:"type:varchar(20);not null" json:"position"`
	JerseyNumber int            `gorm:"type:integer;not null;uniqueIndex:idx_unique_player_jersey_per_team,where:deleted_at IS NULL" json:"jersey_number"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreatePlayerRequest struct {
	TeamID       string  `json:"team_id" binding:"required,uuid"`
	Name         string  `json:"name" binding:"required"`
	HeightCm     float64 `json:"height_cm" binding:"required,gt=0,lt=300"`
	WeightKg     float64 `json:"weight_kg" binding:"required,gt=0,lt=200"`
	Position     string  `json:"position" binding:"required,oneof=forward midfielder defender goalkeeper"`
	JerseyNumber int     `json:"jersey_number" binding:"required,gt=0,lte=99"`
}

type UpdatePlayerRequest struct {
	TeamID       string  `json:"team_id" binding:"required,uuid"`
	Name         string  `json:"name" binding:"required"`
	HeightCm     float64 `json:"height_cm" binding:"required,gt=0,lt=300"`
	WeightKg     float64 `json:"weight_kg" binding:"required,gt=0,lt=200"`
	Position     string  `json:"position" binding:"required,oneof=forward midfielder defender goalkeeper"`
	JerseyNumber int     `json:"jersey_number" binding:"required,gt=0,lte=99"`
}

type GetPlayerRequest struct {
	ID           string
	TeamID       string
	JerseyNumber int
}
