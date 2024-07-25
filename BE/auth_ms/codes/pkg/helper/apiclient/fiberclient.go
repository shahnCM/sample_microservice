package apiclient

import (
	"auth_ms/pkg/dto/apiclient"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func MakeRequest(apiClientRequestDto *apiclient.ApiClientRequestDto) (*apiclient.ApiClientResponseDto, error) {

	method := strings.ToLower(apiClientRequestDto.Method)

	switch method {

	case "get":
		return get(apiClientRequestDto)

	case "post":
		return nil, nil

	case "put":
		return nil, nil

	case "delete":
		return nil, nil

	default:
		return nil, fiber.ErrUnprocessableEntity
	}
}

func checkErrorAndProcessData(agent *fiber.Agent) (*apiclient.ApiClientResponseDto, error) {
	// agent.Debug()
	statusCode, bodyByte, errs := agent.Bytes()
	if len(errs) > 0 {
		var errStr string
		for _, err := range errs {
			errStr += err.Error()
		}

		return nil, fiber.NewError(422, errStr)
	}

	var data map[string]any
	if err := json.Unmarshal(bodyByte, &data); err != nil {
		errStr := "FiberClient: Recieved StatusCode: " + strconv.Itoa(statusCode) + ", Error in response body, unable to unmarshal: " + err.Error()
		return nil, fiber.NewError(422, errStr)
	}

	return &apiclient.ApiClientResponseDto{BodyByte: bodyByte, StatusCode: statusCode, Headers: nil, Body: data}, nil
}

func get(apiClientRequestDto *apiclient.ApiClientRequestDto) (*apiclient.ApiClientResponseDto, error) {
	agent := fiber.Get(apiClientRequestDto.ApiUrl).Timeout(time.Duration(apiClientRequestDto.Timeout) * time.Second)
	return checkErrorAndProcessData(agent)
}

func post() (any, error) {
	return nil, nil
}

func put() (any, error) {
	return nil, nil
}

func delete() (any, error) {
	return nil, nil
}
