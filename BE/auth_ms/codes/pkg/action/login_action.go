package action

import (
	"auth_ms/pkg/dto"
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/helper/safeasync"
	"auth_ms/pkg/provider/database/mariadb10"
	"auth_ms/pkg/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Login(userLoginReqP *request.UserLoginDto) (*dto.UserTokenDataDto, *fiber.Error) {
	// declaring service response and error
	var err error

	// Verify Username
	userService := service.NewUserService(nil)
	userModelP, err := userService.GetUserByUsername(userLoginReqP)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

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
	serviceResponseP, err := service.IssueJwtWithRefreshToken(userModelP.Id, userModelP.Role, hashedUlid, nil)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	tokenDataP := serviceResponseP.Data.(*dto.TokenDataDto)

	defer safeasync.Run(func() {

		tx := mariadb10.GetMariaDb10().Begin()
		if err = tx.Error; err != nil {
			// return nil, fiber.ErrInternalServerError
			return
		}

		err = func(tx *gorm.DB) error {
			// Create a New Associated Session
			sessionService := service.NewSessionService(tx)
			sessionModelP, err := sessionService.StartSession(&userModelP.Id, ulidP, tokenDataP.Jwt.TokenExp, tokenDataP.Refresh.TokenExp)
			if err != nil {
				return err
			}

			// Update user with new session id and new token id
			userService = service.NewUserService(tx)
			_, err = userService.StartUserActiveSessionAndToken(userModelP, &sessionModelP.Id, ulidP)
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
		}(tx)

		if err != nil {
			if err = tx.Rollback().Error; err != nil {
				log.Println(err.Error())
			}
			// return nil, fiber.ErrInternalServerError
			return
		}

		if err = tx.Commit().Error; err != nil {
			log.Println(err.Error())
			// return nil, fiber.ErrInternalServerError
			return
		}
	})

	return &dto.UserTokenDataDto{
		TokenDataDto: tokenDataP,
		User: &map[string]any{
			"username": userModelP.Username,
			"role":     userModelP.Role,
		},
	}, nil
}
