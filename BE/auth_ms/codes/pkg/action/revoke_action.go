package action

import (
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/helper/safeasync"
	"auth_ms/pkg/provider/database/mariadb10"
	"auth_ms/pkg/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Revoke(jwtToken *string) *fiber.Error {

	// Verify token and get claims
	responseP, err := service.VerifyJWT(jwtToken)
	if err != nil {
		return fiber.NewError(401, "Invalid Jwt Token")
	}
	claims := responseP.Data.(*service.Claims)

	defer safeasync.Run(func() {

		// Begin Transaction
		tx := mariadb10.GetMariaDb10().Begin()
		if err := tx.Error; err != nil {
			return
		}

		err := func() error {

			// Fetch user by user_id
			userService := service.NewUserService(tx)
			userModelP, err := userService.GetUserById(&claims.UserId, false)
			if err != nil {
				return fiber.ErrUnauthorized
			}

			if userModelP.SessionTokenTraceId == nil || !common.CompareHash(claims.TokenId, userModelP.SessionTokenTraceId) {
				return fiber.ErrUnauthorized
			}

			// Update user's active session to NULL Only if its associated with the requested token
			_, err = userService.EndUserActiveSessionAndToken(userModelP)
			if err != nil {
				return err
			}

			// End session associated with the token
			sessionService := service.NewSessionService(tx)
			_, err = sessionService.EndSession(userModelP.LastSessionId)
			if err != nil {
				return err
			}

			return nil
		}()

		if err != nil {
			log.Println("ERROR: ", err.Error())
			if errRollback := tx.Rollback().Error; errRollback != nil {
				log.Println("ERROR Rollback: ", errRollback)
			}
			return
		}

		if errCommit := tx.Commit().Error; errCommit != nil {
			log.Println("ERROR Commit: ", errCommit)
			return
		}

	})

	return nil
}
