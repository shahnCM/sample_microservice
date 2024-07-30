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
)

func Login(userLoginReqP *request.UserLoginDto) (*dto.UserTokenDataDto, *fiber.Error) {

	ulidPChan := make(chan *string)
	defer close(ulidPChan)

	hashedUlidPChan := make(chan *string)
	defer close(hashedUlidPChan)

	passMatchChan := make(chan bool)
	defer close(passMatchChan)

	// Generate ULID for Token
	safeasync.Run(func() {
		ulidP, err := common.GenerateULID()
		if err != nil {
			ulidPChan <- nil
			hashedUlidPChan <- nil
			return
		}
		ulidPChan <- ulidP

		// Generate Hash for Ulid
		safeasync.Run(func() {
			hashedUlidP, err := common.GenerateHash(ulidP)
			if err != nil {
				hashedUlidPChan <- nil
				return
			}
			hashedUlidPChan <- hashedUlidP
		})
	})

	// Verify Username
	userService := service.NewUserService(nil)
	userModelP, err := userService.GetUserByUsername(userLoginReqP)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	// compare pass asynchronously
	safeasync.Run(func() {
		if !common.CompareHash(&userModelP.Password, &userLoginReqP.Password) {
			passMatchChan <- false
			return
		}
		passMatchChan <- true
	})

	ulidP := <-ulidPChan
	if ulidP == nil {
		return nil, fiber.ErrUnauthorized
	}
	hashedUlidP := <-hashedUlidPChan
	if ulidP == nil {
		return nil, fiber.ErrUnauthorized
	}

	// Set claims & Generate a new JWT token and Associated Refresh Token
	serviceResponseP, err := service.IssueJwtWithRefreshToken(userModelP.Id, userModelP.Role, hashedUlidP, nil)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	tokenDataP := serviceResponseP.Data.(*dto.TokenDataDto)

	if passMatched := <-passMatchChan; !passMatched {
		return nil, fiber.ErrUnauthorized
	}

	safeasync.Run(func() {

		tx := mariadb10.GetMariaDb10().Begin()
		if err = tx.Error; err != nil {
			return // return nil, fiber.ErrInternalServerError
		}

		err = func() error {
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
		}()

		if err != nil {
			if err = tx.Rollback().Error; err != nil {
				log.Println(err.Error())
			}
			return // nil, fiber.ErrInternalServerError
		}

		if err = tx.Commit().Error; err != nil {
			log.Println(err.Error())
			return // nil, fiber.ErrInternalServerError
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
