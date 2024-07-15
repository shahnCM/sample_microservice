package apiclient

type ApiClientRequestDto struct {
	ApiUrl      string `json:"api_url"`
	Method      string `json:"method"`
	QueryParams string `json:"query_params"`
	Timeout     int    `json:"timeout"`
	Headers     any    `json:"headers"`
	Body        any    `json:"body"`
}

type ApiClientResponseDto struct {
	BodyByte   []byte         `json:"byte_data"`
	StatusCode int            `json:"status_code"`
	Headers    any            `json:"headers"`
	Body       map[string]any `json:"body"`
}
