package dto

type TokenDataDto struct {
	Jwt     *TokenDto `json:"jwt"`
	Refresh *TokenDto `json:"refresh"`
}

type TokenDto struct {
	Type     string  `json:"type"`
	Token    *string `json:"token"`
	TokenExp *int64  `json:"expires_at"`
}
