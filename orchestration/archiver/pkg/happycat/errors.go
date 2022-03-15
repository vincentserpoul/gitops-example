package happycat

import (
	"archiver/pkg/handler"
	"errors"
	"fmt"
	"net/http"
)

type NoHappyCatFactFoundError struct{}

func (e NoHappyCatFactFoundError) Error() string {
	return "no happy cat fact found"
}

type WrongIDFormatError struct {
	ID  string
	Err error
}

func (e WrongIDFormatError) Error() string {
	return fmt.Sprintf("%s is not a valid id: %v", e.ID, e.Err)
}

type WrongBodyFormatError struct {
	Err error
}

func (e WrongBodyFormatError) Error() string {
	return fmt.Sprintf("submitted body is invalid: %v", e.Err)
}

func newErrorResponse(err error) *handler.ErrorResponse {
	switch {
	case errors.As(err, &NoHappyCatFactFoundError{}):
		return &handler.ErrorResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorUser:      handler.ErrorUser("no happy cat facts have been found"),
		}

	case errors.As(err, &WrongIDFormatError{}):
		return &handler.ErrorResponse{
			Err:            err,
			HTTPStatusCode: http.StatusBadRequest,
			ErrorUser:      handler.ErrorUser(err.Error()),
		}

	case errors.As(err, &WrongBodyFormatError{}):
		return &handler.ErrorResponse{
			Err:            err,
			HTTPStatusCode: http.StatusBadRequest,
			ErrorUser: handler.ErrorUser(
				`
wrong format for the happy cat fact. 
should be
{
	'id': '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
	'fact': 'fact'
}`,
			),
		}

	default:
		return &handler.ErrorResponse{
			Err:            err,
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorUser:      handler.ErrorUser("an internal error has occurred"),
		}
	}
}
