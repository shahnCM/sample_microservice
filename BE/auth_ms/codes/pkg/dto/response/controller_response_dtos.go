package response

type SuccessResponseDto struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
	Meta       any    `json:"meta"`
	Links      any    `json:"links"`
}

type ErrorResponseDto struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Message    any    `json:"message"`
}

type ValidationErrorResponseDto struct {
	Status           string                        `json:"status"`
	StatusCode       int                           `json:"status_code"`
	Message          string                        `json:"message"`
	Errors           *map[string]*[]string         `json:"errors"`
	ValidationErrors *[]*ValidationErrorElementDto `json:"validation_errors"`
}

type ValidationErrorElementDto struct {
	Msg                  string `json:"msg"`
	FailedField          string `json:"failed_field"`
	FailedFieldInputType string `json:"failed_field_input_type"`
	FailedFieldReason    string `json:"failed_field_reason"`
	// Tag         string `json:"tag"`
	// Value       string `json:"value"`
}
