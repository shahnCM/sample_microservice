package request

type UserRegistrationDto struct {
	Email           string `json:"email" validate:"required"`
	Username        string `json:"username" validate:"required,min=3,max=32"`
	Password        string `json:"password" validate:"required,min=6,max=16"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}
