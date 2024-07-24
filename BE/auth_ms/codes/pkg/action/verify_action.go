package action

import (
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/model"
	"auth_ms/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func Verify(jwtToken *string) (any, *fiber.Error) {

	// Verify token and get claims
	responseP, err := service.VerifyJWT(jwtToken)
	if err != nil {
		return nil, fiber.NewError(401, "Invalid Refresh/Jwt token: Malformed")
	}
	if responseP.StatusCode != 200 {
		return nil, fiber.NewError(401, "Invalid Refresh/Jwt token: Expired")
	}
	claims := responseP.Data.(*service.Claims)

	// Fetch user by user_id from claims
	userService := service.NewUserService(nil)
	resultP, err := userService.GetUserById(&claims.UserId)
	if err != nil {
		return nil, fiber.NewError(401, "Invalid Refresh/Jwt token")
	}
	userModelP := resultP.(*model.User)

	// Check if user's active token_id matches with the claim's token_id
	if userModelP.SessionTokenTraceId == nil || !common.CompareHash(claims.TokenId, userModelP.SessionTokenTraceId) {
		return nil, fiber.NewError(401, "Invalid Refresh/Jwt token")
	}

	return &map[string]any{
		"username": userModelP.Username,
		"role":     userModelP.Role,
	}, nil
}
