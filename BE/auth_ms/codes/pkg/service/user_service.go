package service

import (
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
	"log"

	"gorm.io/gorm"
)

type UserService interface {
	GetUserByIdFast(userIdP *uint) (*model.User, error)
	GetUserById(userIdP *uint, lockForUpdate bool) (*model.User, error)
	GetUserByUsername(userP *request.UserLoginDto) (*model.User, error)
	RegisterUser(userP *request.UserRegistrationDto) (any, error)
	StartUserActiveSessionAndToken(userModelP *model.User, sessionId *uint, sessionTokenTraceId *string) (any, error)
	UpdateUserActiveToken(userModelP *model.User, sessionTokenTraceId *string) (any, error)
	EndUserActiveSessionAndToken(userModelP *model.User) (any, error)
}

func NewUserService(newTx *gorm.DB) UserService {
	if newTx != nil {
		return &baseService{tx: newTx}
	}

	return &baseService{tx: nil}
}

func (s *baseService) GetUserByIdFast(userIdP *uint) (*model.User, error) {
	userRepo := repository.NewUserRepository(s.tx)
	userModelP, err := userRepo.FindUserByIdFast(userIdP)
	if err != nil {
		return nil, err
	}

	return userModelP, nil
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
		log.Println(err)
		return nil, err
	}

	return userModelP, nil
}

func (s *baseService) GetUserByUsername(userP *request.UserLoginDto) (*model.User, error) {
	userRepo := repository.NewUserRepository(s.tx)
	user, err := userRepo.FindUser(userP.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *baseService) RegisterUser(userP *request.UserRegistrationDto) (any, error) {
	userModelP := &model.User{
		Username:            userP.Username,
		Email:               userP.Email,
		Password:            userP.Password,
		Role:                enum.STANDARD,
		LastSessionId:       nil,
		SessionTokenTraceId: nil,
	}

	userRepo := repository.NewUserRepository(s.tx)
	err := userRepo.CreateUser(userModelP)
	if err != nil {
		return nil, err
	}

	return userModelP, nil
}

func (s *baseService) StartUserActiveSessionAndToken(userModelP *model.User, sessionId *uint, sessionTokenTraceId *string) (any, error) {
	userRepo := repository.NewUserRepository(s.tx)
	userModelP.LastSessionId = sessionId
	userModelP.SessionTokenTraceId = sessionTokenTraceId
	err := userRepo.UpdateUser(userModelP)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *baseService) UpdateUserActiveToken(userModelP *model.User, sessionTokenTraceId *string) (any, error) {
	userModelP.SessionTokenTraceId = sessionTokenTraceId
	userRepo := repository.NewUserRepository(s.tx)
	err := userRepo.UpdateUser(userModelP)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *baseService) EndUserActiveSessionAndToken(userModelP *model.User) (any, error) {
	userRepo := repository.NewUserRepository(s.tx)
	userModelP.LastSessionId = nil
	userModelP.SessionTokenTraceId = nil
	err := userRepo.UpdateUser(userModelP)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
