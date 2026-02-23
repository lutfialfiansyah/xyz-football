package usecase

import (
	"context"

	"xyz-football-api/internal/modules/admin"
	"xyz-football-api/internal/modules/admin/repository/postgres"
)

type AdminUsecase interface {
	Login(ctx context.Context, req admin.LoginRequest) (admin.LoginResponse, error)
}

type adminUsecase struct {
	adminRepo     postgres.AdminRepository
	jwtSecret     string
	jwtExpiration int
}

func NewAdminUsecase(repo postgres.AdminRepository, secret string, expirationHours int) AdminUsecase {
	return &adminUsecase{
		adminRepo:     repo,
		jwtSecret:     secret,
		jwtExpiration: expirationHours,
	}
}
