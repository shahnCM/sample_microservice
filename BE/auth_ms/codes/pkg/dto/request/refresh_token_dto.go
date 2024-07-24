package request

type RefreshTokenDto struct {
	Type  string  `json:"type" validate:"required"`
	Token *string `json:"token" validate:"required"`
}
