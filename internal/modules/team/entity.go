package team

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex:idx_unique_teams_name,where:deleted_at IS NULL;not null" json:"name"`
	LogoURL     string         `gorm:"type:text" json:"logo_url"`
	FoundedYear int            `gorm:"type:integer" json:"founded_year"`
	Address     string         `gorm:"type:text" json:"address"`
	City        string         `gorm:"type:varchar(100);index:idx_teams_city,where:deleted_at IS NULL;not null" json:"city"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateTeamRequest struct {
	Name        string                `form:"name" binding:"required"`
	Logo        *multipart.FileHeader `form:"logo" binding:"required"`
	FoundedYear int                   `form:"founded_year" binding:"required,gt=1800"`
	Address     string                `form:"address"`
	City        string                `form:"city" binding:"required"`
}

type UpdateTeamRequest struct {
	Name        string                `form:"name" binding:"required"`
	Logo        *multipart.FileHeader `form:"logo"`
	FoundedYear int                   `form:"founded_year" binding:"required,gt=1800"`
	Address     string                `form:"address"`
	City        string                `form:"city" binding:"required"`
}

type GetTeamRequest struct {
	ID   string
	Name string
}
