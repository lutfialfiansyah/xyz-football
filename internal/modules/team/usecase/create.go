package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"xyz-football-api/internal/modules/team"
	"xyz-football-api/internal/pkg/apperror"
	"xyz-football-api/internal/pkg/utils"

	"gorm.io/gorm"
)

func (u *teamUsecase) CreateTeam(ctx context.Context, req team.CreateTeamRequest) (*team.Team, error) {
	if req.Logo == nil {
		return nil, apperror.New(http.StatusBadRequest, "logo file is required", "file logo wajib diisi")
	}

	// Check name uniqueness beforehand
	existingTeam, err := u.teamRepo.Get(ctx, &team.GetTeamRequest{Name: req.Name})
	if err == nil && existingTeam != nil {
		return nil, apperror.New(http.StatusBadRequest, "team name already exists", "nama tim sudah ada")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	ext := filepath.Ext(req.Logo.Filename)
	slug := utils.Slugify(req.Name)
	filename := fmt.Sprintf("%s_%d%s", slug, time.Now().Unix(), ext)

	// Upload using the injected storage provider
	logoURL, err := u.storage.Upload(req.Logo, "uploads/team/logos", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to upload logo: %w", err)
	}

	newTeam := &team.Team{
		Name:        req.Name,
		LogoURL:     logoURL,
		FoundedYear: req.FoundedYear,
		Address:     req.Address,
		City:        req.City,
	}

	err = u.teamRepo.Create(ctx, newTeam)
	if err != nil {
		return nil, err
	}

	// Format LogoURL to absolute URL before returning
	newTeam.LogoURL = u.storage.GetURL(newTeam.LogoURL)

	return newTeam, nil
}
