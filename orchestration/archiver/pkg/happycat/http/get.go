package http

import (
	handlerhttp "archiver/pkg/handler/http"
	"archiver/pkg/internal/db"
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type NoHappyCatFactFoundError struct{}

func (e NoHappyCatFactFoundError) Error() string {
	return "no happy cat fact found"
}

func getHandler(
	q db.Querier,
	gnup handlerhttp.NamedURLParamsGetter,
) handlerhttp.TypedHandler {
	return func(r *http.Request) (*handlerhttp.Response, *handlerhttp.ErrorResponse) {
		happyCatFactID, errR := gnup(r.Context(), "happyCatFactID")
		if errR != nil {
			return nil, errR
		}

		id, err := uuid.Parse(happyCatFactID)
		if err != nil {
			return nil, handlerhttp.ParsingParamError{
				Name:  "happyCatFactID",
				Value: happyCatFactID,
			}.ToErrorResponse()
		}

		hcf, err := q.GetHappycatFact(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, handlerhttp.NotFoundError{Designation: "happy cat fact"}.ToErrorResponse()
			}

			return nil, handlerhttp.InternalServerError{Err: err}.ToErrorResponse()
		}

		return &handlerhttp.Response{
			Body:           hcf,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}
