package services

import (
	"context"
	"ewallet-ums/helpers"
	"ewallet-ums/internal/interfaces"
	"ewallet-ums/internal/models"
	"time"

	"github.com/pkg/errors"
)

type RefreshTokenService struct {
	UserRepo interfaces.IUserRepository
}

func (s *RefreshTokenService) RefreshToken(ctx context.Context, refreshToken string, tokenClaim helpers.ClaimToken) (models.RefreshTokenResponse, error) {
	resp := models.RefreshTokenResponse{}
	token, err := helpers.GenerateToken(ctx, tokenClaim.UserID, tokenClaim.Username, tokenClaim.Fullname, "token", time.Now())
	if err != nil {
		return resp, errors.Wrap(err, "failed to generate new token")
	}

	err = s.UserRepo.UpdateTokenWByRefreshToken(ctx, token, refreshToken)
	if err != nil {
		return resp, errors.Wrap(err, "failed to update new token")
	}
	resp.Token = token
	return resp, nil
}