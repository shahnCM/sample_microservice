package service

import (
	"auth_ms/pkg/dto"
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
)

func GetToken(tokenIdP *string) (*response.GenericServiceResponseDto, error) {
	tokenRepo := repository.NewTokenRepository()
	token, err := tokenRepo.FindToken(tokenIdP)
	if err != nil {
		return nil, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: token}, nil
}

func StoreToken(userIdP *uint, sessionIdP *uint, ulidP *string, tokenDataP *dto.TokenDataDto) (*response.GenericServiceResponseDto, error) {
	tokenModelP := &model.Token{
		Id:               ulidP,
		UserId:           userIdP,
		SessionId:        sessionIdP,
		TokenStatus:      "fresh",
		JwtExpiresAt:     tokenDataP.Jwt.TokenExp,
		RefreshExpiresAt: tokenDataP.Refresh.TokenExp,
		JwtToken:         tokenDataP.Jwt.Token,
		RefreshToken:     tokenDataP.Refresh.Token,
	}

	tokenRepo := repository.NewTokenRepository()
	err := tokenRepo.SaveToken(tokenModelP)
	if err != nil {
		return nil, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: tokenModelP}, nil
}

func UpdateTokenStatus(tokenIdP *string, tokenStatus string) (*response.GenericServiceResponseDto, error) {
	tokenRepo := repository.NewTokenRepository()
	token, err := tokenRepo.UpdateTokenStatus(tokenIdP, tokenStatus)
	if err != nil {
		return nil, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: token}, nil
}
