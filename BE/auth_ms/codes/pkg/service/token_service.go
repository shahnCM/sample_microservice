package service

import (
	"auth_ms/pkg/dto"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
	"time"

	"gorm.io/gorm"
)

type TokenService interface {
}

func NewTokenService(newTx *gorm.DB) TokenService {
	if tx != nil {
		return &baseService{tx: newTx}
	}

	return &baseService{tx: nil}
}

func (s *baseService) GetToken(tokenIdP *string) (any, error) {
	tokenRepo := repository.NewTokenRepository(s.tx)
	token, err := tokenRepo.FindToken(tokenIdP)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *baseService) StoreToken(userIdP *uint, sessionIdP *uint, ulidP *string, tokenDataP *dto.TokenDataDto) (any, error) {
	tokenModelP := &model.Token{
		Id:               ulidP,
		UserId:           userIdP,
		SessionId:        sessionIdP,
		TokenStatus:      enum.FRESH_TOKEN,
		JwtExpiresAt:     time.Unix(*tokenDataP.Jwt.TokenExp, 0),
		RefreshExpiresAt: time.Unix(*tokenDataP.Refresh.TokenExp, 0),
		JwtToken:         tokenDataP.Jwt.Token,
		RefreshToken:     tokenDataP.Refresh.Token,
	}

	tokenRepo := repository.NewTokenRepository(s.tx)
	err := tokenRepo.SaveToken(tokenModelP)
	if err != nil {
		return nil, err
	}

	return tokenModelP, nil
}

func (s *baseService) UpdateTokenStatus(tokenIdP *string, tokenStatus string) (any, error) {
	tokenRepo := repository.NewTokenRepository(s.tx)
	err := tokenRepo.UpdateTokenStatus(tokenIdP, tokenStatus)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
