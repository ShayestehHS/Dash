package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"Dash/pkg/auth"
	"Dash/pkg/logger"
)

var (
	ErrCredentialsNotValid     = errors.New("credentials are not valid")
	ErrInputCredentialsInvalid = errors.New("input credentials are invalid")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(phoneNumber PhoneNumber, password string) (*Token, error) {
	user, err := s.repo.GetByPhone(phoneNumber)
	if err != nil {
		logger.Error("Login:USER_NOT_FOUND", map[string]interface{}{
			"phone": phoneNumber.String(),
			"error": err.Error(),
		})
		return nil, ErrCredentialsNotValid
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		logger.Error("Login:PASSWORD_MISMATCH", map[string]interface{}{
			"phone":   phoneNumber.String(),
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return nil, ErrInputCredentialsInvalid
	}

	baseClaim := auth.NewBaseClaim(user.ID)

	accessToken, err := auth.GenerateAccessToken(baseClaim)
	if err != nil {
		logger.Error("Login:ACCESS_TOKEN_GEN_FAILED", map[string]interface{}{
			"user_id": user.ID,
			"phone":   phoneNumber.String(),
			"error":   err.Error(),
		})
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken(baseClaim)
	if err != nil {
		logger.Error("Login:REFRESH_TOKEN_GEN_FAILED", map[string]interface{}{
			"user_id": user.ID,
			"phone":   phoneNumber.String(),
			"error":   err.Error(),
		})
		return nil, err
	}

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
