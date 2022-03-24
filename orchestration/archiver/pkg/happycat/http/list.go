package http

import (
	handlerhttp "archiver/pkg/handler/http"
	"archiver/pkg/internal/db"
	"database/sql"
	"errors"
	"net/http"
)

// listHandler renders the article from the context
// @Summary Get list of happy cat facts
// @Description listHandler returns a list of cat fact
// @Tags happyCatFact
// @Produce json
// @Router /happycatfact [get]
// @Success 200 {object} []db.HappycatFact
// @Failure 400 {object} handlerhttp.ErrorResponse
// @Failure 404 {object} handlerhttp.ErrorResponse
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
