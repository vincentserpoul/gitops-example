package http

import (
	"context"
	"fmt"
	"net/http"
)

type NamedURLParamsGetter func(ctx context.Context, key string) (string, *ErrorResponse)

type MissingParamError struct {
	Name string
}

func (e MissingParamError) Error() string {
	return fmt.Sprintf("named URL param `%s` is missing", e.Name)
}

func (e MissingParamError) ToErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Error:          e,
		HTTPStatusCode: http.StatusBadRequest,
		ErrorCode:      "missing_param_error",
		ErrorMsg:       e.Error(),
	}
}

type ParsingParamError struct {
	Name  string
	Value string
}

func (e ParsingParamError) Error() string {
	return fmt.Sprintf("can not parse named URL param `%s`: `%s` is invalid", e.Name, e.Value)
}

func (e ParsingParamError) ToErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Error:          e,
		HTTPStatusCode: http.StatusBadRequest,
		ErrorCode:      "parsing_param_error",
		ErrorMsg:       e.Error(),
	}
}
