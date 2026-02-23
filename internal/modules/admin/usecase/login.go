package usecase

import (
	"context"
	"net/http"
	"time"

	"xyz-football-api/internal/modules/admin"
	"xyz-football-api/internal/pkg/apperror"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (u *adminUsecase) Login(ctx context.Context, req admin.LoginRequest) (admin.LoginResponse, error) {
	adm, err := u.adminRepo.Get(ctx, &admin.GetAdminRequest{Username: req.Username})
	if err != nil {
		return admin.LoginResponse{}, apperror.New(http.StatusUnauthorized, "invalid username", "username tidak valid")
	}

	err = bcrypt.CompareHashAndPassword([]byte(adm.PasswordHash), []byte(req.Password))
	if err != nil {
		return admin.LoginResponse{}, apperror.New(http.StatusUnauthorized, "invalid password", "password tidak valid")
	}

	// Generate JWT Token
	expirationTime := time.Now().Add(time.Duration(u.jwtExpiration) * time.Second)
	claims := &admin.Claims{
		ID:       adm.ID,
		Username: adm.Username,
		Role:     adm.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return admin.LoginResponse{}, err
	}

	return admin.LoginResponse{
		Token: tokenString,
		Role:  adm.Role,
	}, nil
}
