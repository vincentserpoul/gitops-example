package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type ErrorResponse struct {
	Error          error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	ErrorCode      string `json:"error_code"`
	ErrorMsg       string `json:"error_msg"`
}

func NewErrorResponse(
	e error,
	hsc int,
	ec string,
	msg string,
) *ErrorResponse {
	return &ErrorResponse{
		Error:          e,
		HTTPStatusCode: hsc,
		ErrorCode:      ec,
		ErrorMsg:       msg,
	}
}

func (her *ErrorResponse) render(log zerolog.Logger, w http.ResponseWriter, r *http.Request) {
	render(
		log,
		r.Header.Get("Accept"),
		her.HTTPStatusCode,
		her,
		w,
	)
}

func (her *ErrorResponse) IsEqual(e1 *ErrorResponse) bool {
	if !errors.Is(e1.Error, her.Error) {
		return false
	}

	if e1.HTTPStatusCode != her.HTTPStatusCode {
		return false
	}

	if e1.ErrorCode != her.ErrorCode {
		return false
	}

	if e1.ErrorMsg != her.ErrorMsg {
		return false
	}

	return true
}

type InternalServerError struct {
	Err error
}

func (e InternalServerError) Error() string {
	return e.Err.Error()
}

func (e InternalServerError) ToErrorResponse() *ErrorResponse {
	return NewErrorResponse(e, http.StatusInternalServerError, "internal_error", "internal error")
}

type NotFoundError struct {
	Designation string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("no corresponding `%s` has been found", e.Designation)
}

func (e NotFoundError) ToErrorResponse() *ErrorResponse {
	return NewErrorResponse(e, http.StatusNotFound, "not_found", e.Error())
}
