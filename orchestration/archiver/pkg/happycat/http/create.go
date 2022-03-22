package http

import (
	handlerhttp "archiver/pkg/handler/http"
	"archiver/pkg/internal/db"
	"net/http"
)

func createHandler(
	q db.Querier,
) handlerhttp.TypedHandler {
	return func(r *http.Request) (*handlerhttp.Response, *handlerhttp.ErrorResponse) {
		var params db.SaveHappycatFactParams

		if err := handlerhttp.DecodeBody(r, &params); err != nil {
			return nil, &handlerhttp.ErrorResponse{
				Error:          err,
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "bad_request",
				ErrorMsg:       "wrong cat fact post",
			}
		}

		if err := q.SaveHappycatFact(r.Context(), params); err != nil {
			return nil, &handlerhttp.ErrorResponse{
				Error:          err,
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorCode:      "db_error",
				ErrorMsg:       "internal error",
			}
		}

		return &handlerhttp.Response{
			Body:           nil,
			HTTPStatusCode: http.StatusCreated,
		}, nil
	}
}
