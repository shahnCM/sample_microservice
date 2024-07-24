package action

import (
	"auth_ms/pkg/dto"
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"
	"auth_ms/pkg/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Login(userLoginReqP *request.UserLoginDto) (*dto.UserTokenDataDto, *fiber.Error) {
	// declaring service response and error
	var err error

	// Verify Username
	userService := service.NewUserService(nil)
	resultP, err := userService.GetUser(userLoginReqP)
	if err != nil {
		log.Println(err)
		return nil, fiber.ErrUnauthorized
	}
	userModelP := resultP.(*model.User)

	// Compare password
	if !common.CompareHash(&userModelP.Password, &userLoginReqP.Password) {
		return nil, fiber.ErrUnauthorized
	}

	// Generate ULID for Token
	ulidP, err := common.GenerateULID()
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	hashedUlid, err := common.GenerateHash(ulidP)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// Set claims & Generate a new JWT token and Associated Refresh Token
	serviceResponseP, err := service.IssueJwtWithRefreshToken(userModelP.Id, userModelP.Role, hashedUlid)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	tokenDataP := serviceResponseP.Data.(*dto.TokenDataDto)

	tx := mariadb10.GetMariaDb10().Begin()
	if err = tx.Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	func() error {
		// Create a New Associated Session
		sessionService := service.NewSessionService(tx)
		resultP, err = sessionService.StoreSession(&userModelP.Id, ulidP, tokenDataP.Jwt.TokenExp, tokenDataP.Refresh.TokenExp)
		if err != nil {
			return err
		}
		sessionModelP := resultP.(*model.Session)

		// Update user with new session id and new token id
		userService = service.NewUserService(tx)
		_, err = userService.UpdateUserActiveSessionAndToken(&userModelP.Id, &sessionModelP.Id, ulidP)
		if err != nil {
			return err
		}

		// End last active session
		if userModelP.LastSessionId != nil {
			_, err = sessionService.EndSession(userModelP.LastSessionId)
			if err != nil {
				return err
			}
		}

		return nil
	}()

	if err != nil {
		if err = tx.Rollback().Error; err != nil {
			log.Println(err.Error())
		}
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &dto.UserTokenDataDto{
		TokenDataDto: tokenDataP,
		User: &map[string]any{
			"username": userModelP.Username,
			"role":     userModelP.Role,
		},
	}, nil
}
