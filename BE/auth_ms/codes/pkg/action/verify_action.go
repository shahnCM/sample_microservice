package action

import (
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func Verify(jwtToken *string) (any, *fiber.Error) {

	// Verify token and get claims
	responseP, err := service.VerifyJWT(jwtToken)
	if err != nil {
		return nil, fiber.NewError(401, "Invalid Jwt token")
	} else if responseP.StatusCode == 401 {
		return nil, fiber.NewError(401, "Expired Jwt token")
	}
	claims := responseP.Data.(*service.Claims)

	// Fetch user by user_id from claims
	userService := service.NewUserService(nil)
	userModelP, err := userService.GetUserByIdFast(&claims.UserId)
	if err != nil ||
		userModelP.SessionTokenTraceId == nil ||
		!common.CompareHash(claims.TokenId, userModelP.SessionTokenTraceId) {
		return nil, fiber.NewError(401, "Invalid Jwt token")
	}

	return &map[string]any{
		"username": userModelP.Username,
		"role":     userModelP.Role,
	}, nil
}
