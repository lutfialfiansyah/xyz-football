package postgres

import (
	"context"
	"xyz-football-api/internal/modules/admin"

	"gorm.io/gorm"
)

type AdminRepository interface {
	Get(ctx context.Context, req *admin.GetAdminRequest) (*admin.Admin, error)
	Create(ctx context.Context, adm *admin.Admin) error
}

type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepo{db: db}
}

func (r *adminRepo) Get(ctx context.Context, req *admin.GetAdminRequest) (*admin.Admin, error) {
	var adm admin.Admin
	query := r.db.WithContext(ctx).Model(&admin.Admin{})

	if req.ID != "" {
		query = query.Where("id = ?", req.ID)
	}

	if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	}

	err := query.First(&adm).Error
	if err != nil {
		return nil, err
	}
	return &adm, nil
}

func (r *adminRepo) Create(ctx context.Context, adm *admin.Admin) error {
	return r.db.WithContext(ctx).Create(adm).Error
}
