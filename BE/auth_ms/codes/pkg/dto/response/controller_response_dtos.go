package response

type SuccessResponseDto struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Data       any    `json:"data"`
	Meta       any    `json:"meta"`
	Links      any    `json:"links"`
}

type ErrorResponseDto struct {
	Status       string `json:"status"`
	StatusCode   int    `json:"status_code"`
	ErrorMessage any    `json:"error_message"`
}

type ValidationErrorResponseDto struct {
	Status           string                       `json:"status"`
	StatusCode       int                          `json:"status_code"`
	ErrorMessage     string                       `json:"error_message"`
	ValidationErrors []*ValidationErrorElementDto `json:"validation_errors"`
}

type ValidationErrorElementDto struct {
	Msg                  string `json:"msg"`
	FailedField          string `json:"failed_field"`
	FailedFieldInputType string `json:"failed_field_input_type"`
	// Tag         string `json:"tag"`
	// Value       string `json:"value"`
}
