package http

import (
	handlerhttp "archiver/pkg/handler/http"
	"archiver/pkg/internal/db"
	"database/sql"
	"errors"
	"net/http"
)

func listHandler(
	q db.Querier,
) handlerhttp.TypedHandler {
	return func(r *http.Request) (*handlerhttp.Response, *handlerhttp.ErrorResponse) {
		hcfs, err := q.ListHappycatFacts(r.Context())
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
			Body:           hcfs,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}
