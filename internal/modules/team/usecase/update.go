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

func (u *teamUsecase) UpdateTeam(ctx context.Context, id string, req team.UpdateTeamRequest) (*team.Team, error) {
	existingTeam, err := u.teamRepo.Get(ctx, &team.GetTeamRequest{ID: id})
	if err != nil {
		return nil, err
	}

	// Check name uniqueness if name is changed
	if req.Name != existingTeam.Name {
		existingNameTeam, err := u.teamRepo.Get(ctx, &team.GetTeamRequest{Name: req.Name})
		if err == nil && existingNameTeam != nil {
			return nil, apperror.New(http.StatusBadRequest, "team name already exists", "nama tim sudah ada")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	logoURL := existingTeam.LogoURL

	// Handle File Upload if Logo is provided
	if req.Logo != nil {
		ext := filepath.Ext(req.Logo.Filename)
		slug := utils.Slugify(req.Name)
		filename := fmt.Sprintf("%s_%d%s", slug, time.Now().Unix(), ext)

		// Upload using the injected storage provider
		path, err := u.storage.Upload(req.Logo, "uploads/team/logos", filename)
		if err != nil {
			return nil, fmt.Errorf("failed to upload logo: %w", err)
		}
		logoURL = path
	}

	existingTeam.Name = req.Name
	existingTeam.LogoURL = logoURL
	existingTeam.FoundedYear = req.FoundedYear
	existingTeam.Address = req.Address
	existingTeam.City = req.City

	err = u.teamRepo.Update(ctx, existingTeam)
	if err != nil {
		return nil, err
	}

	// Format LogoURL to absolute URL before returning
	existingTeam.LogoURL = u.storage.GetURL(existingTeam.LogoURL)

	return existingTeam, nil
}
