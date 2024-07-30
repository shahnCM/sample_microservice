package action

import (
	"auth_ms/pkg/dto"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/helper/safeasync"
	"auth_ms/pkg/provider/database/mariadb10"
	"auth_ms/pkg/service"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Refresh(jwtToken *string, refreshToken *string) (any, *fiber.Error) {

	// declaring service response and error
	var wg sync.WaitGroup
	wg.Add(2)

	statusResults := make(chan int, 2)
	claimsResults := make(chan *service.Claims, 2)

	// verify jwt token concurrently
	safeasync.Run(func() {
		defer wg.Done()
		responseP, _ := service.VerifyJWT(jwtToken)
		if responseP != nil {
			statusResults <- (responseP.StatusCode - 1)
			claimsResults <- responseP.Data.(*service.Claims)
		} else {
			statusResults <- -1
			claimsResults <- &service.Claims{}
		}
	})

	// verify refresh token concurrently
	safeasync.Run(func() {
		defer wg.Done()
		responseP, _ := service.VerifyRefreshToken(refreshToken)
		if responseP != nil {
			statusResults <- (responseP.StatusCode + 9)
			claimsResults <- responseP.Data.(*service.Claims)
		} else {
			statusResults <- -1
			claimsResults <- &service.Claims{}
		}
	})

	wg.Wait()

	close(statusResults)
	close(claimsResults)

	totalStatusResults := 0
	for result := range statusResults {
		totalStatusResults += result
	}

	if totalStatusResults-209 != 400 {
		return nil, fiber.NewError(422, "Can't Refresh: Jwt token hasn't been expired yet")
	}

	/**
	 * total should be 609 indecates that we received
	 * 401 - 1 = 400 from Jwt Verification
	 * 200 + 9 = 209 from Refresh Verification
	 */

	var claims *service.Claims
	var claimsArr []*service.Claims
	for result := range claimsResults {
		claimsArr = append(claimsArr, result)
	}
	// Checking pair integrity of refresh and jwt token
	if *claimsArr[0].TokenId != *claimsArr[1].TokenId {
		return nil, fiber.NewError(422, "Invalid Refresh/Jwt token")
	}
	if claimsArr[0].Type == enum.JWT_TOKEN {
		claims = claimsArr[0]
	} else {
		claims = claimsArr[1]
	}

	// Begin Transaction
	tx := mariadb10.GetMariaDb10().Begin()
	if err := tx.Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	tokenDataP, err := func(tx *gorm.DB) (any, *fiber.Error) {

		// Getting User and Locking for Update
		userService := service.NewUserService(tx)
		userModelP, err := userService.GetUserById(&claims.UserId, true)
		if err != nil {
			return nil, fiber.NewError(404, "Invalid Refresh/Jwt token: User not found")
		}

		// Check if user's active token_id exists as there's any active session running
		if userModelP.SessionTokenTraceId == nil {
			return nil, fiber.NewError(422, "No active session to refresh")
		}

		// Check if user's active token_id matches with the claim's token_id
		if !common.CompareHash(claims.TokenId, userModelP.SessionTokenTraceId) {
			return nil, fiber.NewError(422, "Invalid Refresh/Jwt token")
		}

		// Hashing sessionTokenTraceId
		hashedUlid, err := common.GenerateHash(userModelP.SessionTokenTraceId)
		if err != nil {
			return nil, fiber.NewError(422, "Hash error")
		}

		// Getting Session and Locking for Update
		sessionService := service.NewSessionService(tx)
		sessionModelP, err := sessionService.GetSession(userModelP.LastSessionId, true)
		if err != nil {
			return nil, fiber.NewError(404, "Invalid Refresh/Jwt token: User session not found")
		}

		// Compare current session time with claims session time.
		// if current session time is greater than claims session time than pass current session time as exp for new token
		// otherwise pass nil

		var tokenExp *int64
		dbSessionEndTimeUnix := sessionModelP.EndsAt.Unix()
		if dbSessionEndTimeUnix > *claims.Exp {
			tokenExp = &dbSessionEndTimeUnix
		}

		// Set claims & Generate a new JWT token and Associated Refresh Token
		responseP, err := service.IssueJwtWithRefreshToken(userModelP.Id, userModelP.Role, hashedUlid, tokenExp)
		if err != nil {
			return nil, fiber.NewError(500, "Internal server error: Can't issue token at this moment")
		}
		tokenDataP := responseP.Data.(*dto.TokenDataDto)

		// Update active Associated Session's SessionTokenTraceId, jwtExpiresAt, refreshExpiresAt, refreshCount
		if err = sessionService.RefreshSession(
			sessionModelP, userModelP.SessionTokenTraceId,
			tokenDataP.Jwt.TokenExp, tokenDataP.Refresh.TokenExp,
			&sessionModelP.RefreshCount); err != nil {
			return nil, fiber.NewError(500, "Internal server error: Can't issue token at this moment")
		}

		return tokenDataP, nil
	}(tx)

	if err != nil {
		log.Println("ERROR: ", err.Error())
		if errRollback := tx.Rollback().Error; errRollback != nil {
			log.Println("ERROR Rollback: ", errRollback)
		}
		return nil, err
	}

	if errCommit := tx.Commit().Error; errCommit != nil {
		log.Println("ERROR Commit: ", errCommit)
		return nil, fiber.ErrInternalServerError
	}

	return tokenDataP, nil
}
