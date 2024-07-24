package common

import (
	"auth_ms/pkg/dto/response"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"

	crypto_rand "crypto/rand"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/strcase"
)

func HandleResponse(ctx *fiber.Ctx, res *response.GenericServiceResponseDto) error {
	if res.StatusCode < 300 || res.StatusCode == 304 {
		return SuccessResponse(ctx, res.StatusCode, res.Data, res.Meta, res.Links)
	}

	return ErrorResponse(ctx, res.StatusCode, res)
}

func SuccessResponse(ctx *fiber.Ctx, code int, data any, meta any, links any) error {
	return ctx.Status(code).JSON(&response.SuccessResponseDto{
		Status:     "Success",
		StatusCode: code,
		Data:       data,
		Meta:       meta,
		Links:      links,
	})
}

func ErrorResponse(ctx *fiber.Ctx, code int, message any) error {
	return ctx.Status(code).JSON(&response.ErrorResponseDto{
		Status:       "Error",
		StatusCode:   code,
		ErrorMessage: message,
	})
}

func ValidationErrorResponse(ctx *fiber.Ctx, validationErrorElements []*response.ValidationErrorElementDto) error {
	return ctx.Status(422).JSON(&response.ValidationErrorResponseDto{
		Status:           "Error",
		StatusCode:       422,
		ErrorMessage:     "Validation Error",
		ValidationErrors: validationErrorElements,
	})
}

func ValidateRequest(data any) []*response.ValidationErrorElementDto {
	var validate = validator.New()
	var errors []*response.ValidationErrorElementDto

	if err := validate.Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			var element response.ValidationErrorElementDto

			value := err.Param()
			if value == "" {
				value = "missing field or invalid input type"
			}

			failedField := strcase.ToSnake(err.StructField())
			failedFieldInputType := err.Type().String()
			failedFieldReason := err.Tag()

			element.Msg = fmt.Sprintf(`%s, %s: %s`, failedField, err.Tag(), value)
			element.FailedField = failedField
			element.FailedFieldReason = failedFieldReason
			element.FailedFieldInputType = failedFieldInputType

			errors = append(errors, &element)
		}
	}

	return errors
}

func ValidateRequestWithValidationErrorResponse(ctx *fiber.Ctx, data any) error {
	if err := ValidateRequest(data); err != nil {
		return ValidationErrorResponse(ctx, err)
	}

	return nil
}

func ParseRequestBody(ctx *fiber.Ctx, dto any) (any, error) {

	bodyParserError := ctx.BodyParser(dto)

	if validationError := ValidateRequest(dto); validationError != nil {
		return validationError, ValidationErrorResponse(ctx, validationError)
	}

	if bodyParserError != nil {
		log.Println(bodyParserError)
		return bodyParserError, ErrorResponse(ctx, 422, "Invalid Or Empty input: "+bodyParserError.Error())
	}

	return nil, nil
}

func ParseRouteParams(ctx *fiber.Ctx, param string) (string, map[string]string) {
	if param != "" {
		return ctx.Params(param), nil
	}
	return "", ctx.AllParams()
}

func ParseQueryParams(ctx *fiber.Ctx, param string, defaultValue string) (string, map[string]string) {
	if param != "" {
		return ctx.Query(param, defaultValue), nil
	}
	return "", ctx.Queries()
}

func ParseHeader(ctx *fiber.Ctx, param string, defaultValue string) (string, map[string][]string) {
	if param != "" {
		return ctx.Get(param, defaultValue), nil
	}
	return "", ctx.GetReqHeaders()
}

func ProcessPageAndSize(ctx *fiber.Ctx) (*int, *int, *int) {
	pageStrVal, _ := ParseQueryParams(ctx, "page", "1")
	sizeStrVal, _ := ParseQueryParams(ctx, "per_page", "10")

	pageIntVal, err := strconv.Atoi(pageStrVal)
	if err != nil || pageIntVal < 1 {
		pageIntVal = 1
	}

	sizeIntVal, err := strconv.Atoi(sizeStrVal)
	if err != nil || sizeIntVal < 1 {
		sizeIntVal = 10
	}

	offset := sizeIntVal * (pageIntVal - 1)

	return &pageIntVal, &sizeIntVal, &offset
}

func CurrentUrl(ctx *fiber.Ctx) string {
	return ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
}

func PaginateResponse(searchKey string, currentUrl string, currentPage int, size int, data any, total int) *response.GenericServiceResponseDto {
	total, from := total, (size*(currentPage-1))+1
	to, prevPage, nextPage, lastPage := from+size-1, currentPage-1, currentPage+1, int(math.Ceil(float64(total)/float64(size)))

	if to > total {
		to = total
	}

	if prevPage < 1 {
		prevPage = 1
	}

	if nextPage > lastPage {
		nextPage = lastPage
	}

	if total == 0 {
		to = 0
		from = 0
		nextPage = 0
	}

	if searchKey != "" {
		currentUrl = currentUrl + "?key=" + searchKey
	}

	var metaMap, linksMap map[string]any
	meta := fmt.Sprintf(`{"path":"%s","current_page":%d,"from":%d,"last_page":%d,"limit":%d,"to":%d,"total":%d}`, currentUrl, currentPage, from, lastPage, size, to, total)
	links := fmt.Sprintf(`{"first":"%s&page=%d&per_page=%d","last":"%s&page=%d&per_page=%d","next":"%s&page=%d&per_page=%d","prev":"%s&page=%d&per_page=%d"}`, currentUrl, 1, size, currentUrl, lastPage, size, currentUrl, nextPage, size, currentUrl, prevPage, size)
	json.Unmarshal([]byte(meta), &metaMap)
	json.Unmarshal([]byte(links), &linksMap)

	return &response.GenericServiceResponseDto{
		StatusCode: 200,
		Data:       data,
		Meta:       metaMap,
		Links:      linksMap,
	}
}

func GenerateULID() (*string, error) {
	timestamp := time.Now().UTC()
	entropy := ulid.Monotonic(crypto_rand.Reader, 0)
	ulid, err := ulid.New(ulid.Timestamp(timestamp), entropy)
	if err != nil {
		log.Printf("failed to generate ULID: %v", err)
		return nil, err
	}

	ulidStr := ulid.String()
	return &ulidStr, nil
}

func GenerateHash(str *string) (*string, error) {
	if str == nil {
		return nil, errors.New("input string cannot be nil")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*str), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedStr := string(hashedBytes)
	return &hashedStr, nil
}

func CompareHash(hashedStr, plainStr *string) bool {
	if hashedStr == nil || plainStr == nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(*hashedStr), []byte(*plainStr))
	return err == nil
}
