package response

type GenericServiceResponseDto struct {
	StatusCode int `json:"status_code"`
	Data       any `json:"data"`
	Meta       any `json:"meta"`
	Links      any `json:"links"`
}
