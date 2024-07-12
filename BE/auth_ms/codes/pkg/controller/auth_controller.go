package controller

import (
	"auth_ms/pkg/helper/common"

	"github.com/gofiber/fiber/v2"
)

func Fresh(ctx *fiber.Ctx) error {
	// Verify Username & Password
	// Generate a new JWT token and Associated Refresh Token
	// Create a New Associated Session
	// return JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, "auth - Fresh", nil, nil)
}

func Refresh(ctx *fiber.Ctx) error {
	// Verify token and get claims
	// Check if the token ulid exists on db and its status is fresh
	// Update this token's status to refreshed
	// Generate a new JWT token and Associated Refresh Token
	// Associate it with the current Running session
	// return JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, "auth - Refresh", nil, nil)
}

func Verify(ctx *fiber.Ctx) error {
	// Verify token and get claims
	// Check if the token ulid exists on db and its status is fresh
	return common.SuccessResponse(ctx, 200, "auth - Verify", nil, nil)
}

func Revoke(ctx *fiber.Ctx) error {
	// Verify token and get claims
	// Check if the token ulid exists on db and its status is fresh
	// Update this token's status to refreshed
	// Generate a new JWT token and Associated Refresh Token
	// Associate it with the current Running session
	// return JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, "auth - Revoke", nil, nil)
}
