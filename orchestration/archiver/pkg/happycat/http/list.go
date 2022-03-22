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
				return nil, handlerhttp.NotFoundError{Designation: "happy cat fact"}.ToErrorResponse()
			}

			return nil, handlerhttp.InternalServerError{Err: err}.ToErrorResponse()
		}

		return &handlerhttp.Response{
			Body:           hcfs,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}
