package common

import (
	"auth_ms/pkg/dto"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/strcase"
)

func HandleResponse(ctx *fiber.Ctx, res *dto.GenericServiceResponseDto) error {
	if res.StatusCode < 300 || res.StatusCode == 304 {
		return SuccessResponse(ctx, res.StatusCode, res.Data, res.Meta, res.Links)
	}

	return ErrorResponse(ctx, res.StatusCode, res)
}

func SuccessResponse(ctx *fiber.Ctx, code int, data any, meta any, links any) error {
	return ctx.Status(code).JSON(&dto.SuccessResponseDto{
		Status:     "Success",
		StatusCode: code,
		Data:       data,
		Meta:       meta,
		Links:      links,
	})
}

func ErrorResponse(ctx *fiber.Ctx, code int, message any) error {
	return ctx.Status(code).JSON(&dto.ErrorResponseDto{
		Status:       "Error",
		StatusCode:   code,
		ErrorMessage: message,
	})
}

func ValidationErrorResponse(ctx *fiber.Ctx, validationErrorElements []*dto.ValidationErrorElementDto) error {
	return ctx.Status(422).JSON(&dto.ValidationErrorResponseDto{
		Status:           "Error",
		StatusCode:       422,
		ErrorMessage:     "Validation Error",
		ValidationErrors: validationErrorElements,
	})
}

func ValidateRequest(data any) []*dto.ValidationErrorElementDto {
	var validate = validator.New()
	var errors []*dto.ValidationErrorElementDto

	if err := validate.Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			var element dto.ValidationErrorElementDto

			value := err.Param()
			if value == "" {
				value = "missing field or invalid input type"
			}

			failedField := strcase.ToSnake(err.StructField())
			failedFieldInputType := err.Type().String()

			element.Msg = fmt.Sprintf(`%s, %s: %s`, failedField, err.Tag(), value)
			element.FailedField = failedField
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

func ProcessPageAndSize(ctx *fiber.Ctx) (int, int) {
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

	return pageIntVal, sizeIntVal
}

func CurrentUrl(ctx *fiber.Ctx) string {
	return ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
}

func PaginateResponse(searchKey string, forList bool, currentUrl string, currentPage int, size int, data any, total int) *dto.GenericServiceResponseDto {
	currentUrl = currentUrl + "?key=" + searchKey
	total, from := total, (size*(currentPage-1))+1
	to, prevPage, nextPage, lastPage := from+size, currentPage-1, currentPage+1, int(math.Ceil(float64(total)/float64(size)))

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

	var meta string
	var metaMap, linksMap map[string]any

	if forList {
		meta = fmt.Sprintf(`{"path":"%s","current_page":%d,"from":%d,"last_page":%d,"limit":%d,"to":%d,"total":%d}`, currentUrl, currentPage, from, lastPage, size, to, total)
	} else {
		topSearch := `[{"id":1,"key":"internet-pack","value":"Internet Pack","priority":1},{"id":2,"key":"monthly_pack","value":"monthly pack","priority":2}]`
		meta = fmt.Sprintf(`{"top_search":%s, "pagination":{"current_page":%d,"from":%d,"last_page":%d,"limit":%d,"to":%d,"total":%d}}`, topSearch, currentPage, from, lastPage, size, to, total)
	}

	links := fmt.Sprintf(`{"first":"%s&page=%d&per_page=%d","last":"%s&page=%d&per_page=%d","next":"%s&page=%d&per_page=%d","prev":"%s&page=%d&per_page=%d"}`, currentUrl, 1, size, currentUrl, lastPage, size, currentUrl, nextPage, size, currentUrl, prevPage, size)
	json.Unmarshal([]byte(meta), &metaMap)
	json.Unmarshal([]byte(links), &linksMap)

	return &dto.GenericServiceResponseDto{
		StatusCode: 200,
		Data:       data,
		Meta:       metaMap,
		Links:      linksMap,
	}
}
