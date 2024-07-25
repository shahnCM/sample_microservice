package service

import (
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"

	"gorm.io/gorm"
)

type UserService interface {
	GetUserById(userIdP *uint, lockForUpdate bool) (*model.User, error)
	GetUser(userP *request.UserLoginDto) (*model.User, error)
	StoreUser(userP *request.UserRegistrationDto) (any, error)
	UpdateUserActiveToken(userModelP *model.User, tokenIdP *string) (any, error)
	UpdateUserActiveSessionAndToken(userIdP *uint, sessionIdP *uint, tokenIdP *string) (any, error)
}

func NewUserService(newTx *gorm.DB) UserService {
	if newTx != nil {
		return &baseService{tx: newTx}
	}

	return &baseService{tx: nil}
}

func (s *baseService) GetUserById(userIdP *uint, lockForUpdate bool) (*model.User, error) {
	var userModelP *model.User
	var err error

	userRepo := repository.NewUserRepository(s.tx)

	if lockForUpdate {
		userModelP, err = userRepo.FindUserByIdAndLockForUpdate(userIdP)
	} else {
		userModelP, err = userRepo.FindUserById(userIdP)
	}

	if err != nil {
		return nil, err
	}

	return userModelP, nil
}

func (s *baseService) GetUser(userP *request.UserLoginDto) (*model.User, error) {
	userRepo := repository.NewUserRepository(s.tx)
	user, err := userRepo.FindUser(userP.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *baseService) StoreUser(userP *request.UserRegistrationDto) (any, error) {
	userModelP := &model.User{
		Username:            userP.Username,
		Email:               userP.Email,
		Password:            userP.Password,
		Role:                enum.STANDARD,
		LastSessionId:       nil,
		SessionTokenTraceId: nil,
	}

	userRepo := repository.NewUserRepository(s.tx)
	err := userRepo.SaveUser(userModelP)
	if err != nil {
		return nil, err
	}

	return userModelP, nil
}

func (s *baseService) UpdateUserActiveToken(userModelP *model.User, tokenIdP *string) (any, error) {

	userModelP.SessionTokenTraceId = tokenIdP
	userRepo := repository.NewUserRepository(s.tx)
	err := userRepo.SaveUser(userModelP)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *baseService) UpdateUserActiveSessionAndToken(userIdP *uint, sessionIdP *uint, tokenIdP *string) (any, error) {
	userRepo := repository.NewUserRepository(s.tx)
	updatesP := &map[string]any{
		"last_session_id":        sessionIdP,
		"session_token_trace_id": tokenIdP,
	}
	err := userRepo.UpdateUser(userIdP, updatesP)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
