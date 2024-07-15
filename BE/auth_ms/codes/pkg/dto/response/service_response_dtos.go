package response

import "auth_ms/pkg/dto"

type GenericServiceResponseDto struct {
	StatusCode int `json:"status_code"`
	Data       any `json:"data"`
	Meta       any `json:"meta"`
	Links      any `json:"links"`
}

type ElasticsearchGenericResponseDto struct {
	ElasticsearchStatusCode   int `json:"elasticsearch_status_code"`
	ElasticsearchResponseBody any `json:"elasticsearch_response_body"`
}

type ElasticsearchStructuredResponseDto struct {
	ElasticsearchStatusCode   int                        `json:"elasticsearch_status_code"`
	ElasticsearchResponseBody *ElasticsearchResponseBody `json:"elasticsearch_response_body"`
}

type ElasticsearchResponseBody struct {
	Shards   any                `json:"shards"`
	Hits     *ElasticsearchHits `json:"hits"`
	Meta     any                `json:"meta"`
	TimedOut bool               `json:"timed_out"`
	Took     any                `json:"took"`
}

type ElasticsearchHits struct {
	Hits     []*ElasticsearchHitsDocs `json:"hits"`
	MaxScore any                      `json:"max_score"`
	Total    *ElasticsearchTotal      `json:"total"`
}

type ElasticsearchTotal struct {
	Relation string `json:"relation"`
	Value    int    `json:"value"`
}

type ElasticsearchHitsDocs struct {
	Source *dto.Document `json:"_source"`
}
