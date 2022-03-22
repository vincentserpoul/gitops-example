package http

import (
	handlerhttp "archiver/pkg/handler/http"
	"archiver/pkg/internal/db"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

func getHandler(
	q db.Querier,
	gnup handlerhttp.NamedURLParamsGetter,
) handlerhttp.TypedHandler {
	return func(r *http.Request) (*handlerhttp.Response, *handlerhttp.ErrorResponse) {
		happyCatFactID := gnup(r.Context(), "happyCatFactID")

		id, err := uuid.Parse(happyCatFactID)
		if err != nil {
			return nil, &handlerhttp.ErrorResponse{
				Error:          WrongIDFormatError{ID: happyCatFactID, Err: err},
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "bad_request",
				ErrorMsg:       "wrong format id",
			}
		}

		hcf, err := q.GetHappycatFact(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, &handlerhttp.ErrorResponse{
					Error:          NoHappyCatFactFoundError{},
					HTTPStatusCode: http.StatusNotFound,
					ErrorCode:      "db_not_found",
					ErrorMsg:       NoHappyCatFactFoundError{}.Error(),
				}
			}

			return nil, &handlerhttp.ErrorResponse{
				Error:          err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorCode:      "db_error",
				ErrorMsg:       "internal error",
			}
		}

		return &handlerhttp.Response{
			Body:           hcf,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}
