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

		if errR := handlerhttp.BindBody(r, &params); errR != nil {
			return nil, errR
		}

		if err := q.SaveHappycatFact(r.Context(), params); err != nil {
			return nil, handlerhttp.InternalServerError{Err: err}.ToErrorResponse()
		}

		return &handlerhttp.Response{
			Body:           nil,
			HTTPStatusCode: http.StatusCreated,
		}, nil
	}
}
