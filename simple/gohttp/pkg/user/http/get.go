package http

import (
	"gohttp/pkg/user"
	"net/http"

	"github.com/google/uuid"
	"github.com/induzo/httpwrapper"
)

// getHandler renders the article from the context
// @Summary Get user by id
// @Description getHandler returns a single cat fact by id
// @Tags user
// @Produce json
// @Param userID path string true "user id"
// @Router /user/{userID} [get]
// @Success 200 {object} user.User
// @Failure 400 {object} httpwrapper.ErrorResponse
// @Failure 404 {object} httpwrapper.ErrorResponse
func getHandler(
	us *user.Storage,
	gnup httpwrapper.NamedURLParamsGetter,
) httpwrapper.TypedHandler {
	return func(r *http.Request) (*httpwrapper.Response, *httpwrapper.ErrorResponse) {
		userID, errR := gnup(r.Context(), "userID")
		if errR != nil {
			return nil, errR
		}

		id, err := uuid.Parse(userID)
		if err != nil {
			return nil, httpwrapper.ParsingParamError{
				Name:  "userID",
				Value: userID,
			}.ToErrorResponse()
		}

		user, err := us.Get(id)
		if err != nil {
			return nil, httpwrapper.NotFoundError{Designation: "user"}.ToErrorResponse()
		}

		return &httpwrapper.Response{
			Body:           user,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}
