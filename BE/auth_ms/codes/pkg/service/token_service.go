package service

import (
	"auth_ms/pkg/dto"
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
	"time"
)

func GetToken(tokenIdP *string) (*response.GenericServiceResponseDto, error) {
	tokenRepo := repository.NewTokenRepository()
	token, err := tokenRepo.FindToken(tokenIdP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 404, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: token}, nil
}

func StoreToken(userIdP *uint, sessionIdP *uint, ulidP *string, tokenDataP *dto.TokenDataDto) (*response.GenericServiceResponseDto, error) {
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

	tokenRepo := repository.NewTokenRepository()
	err := tokenRepo.SaveToken(tokenModelP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 422, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 201, Data: tokenModelP}, nil
}

func UpdateTokenStatus(tokenIdP *string, tokenStatus string) (*response.GenericServiceResponseDto, error) {
	tokenRepo := repository.NewTokenRepository()
	err := tokenRepo.UpdateTokenStatus(tokenIdP, tokenStatus)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 422, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, nil
}
