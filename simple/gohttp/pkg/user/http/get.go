package http

import (
	handlerhttp "gohttp/pkg/handler/http"
	"gohttp/pkg/user"
	"net/http"

	"github.com/google/uuid"
)

// getHandler renders the article from the context
// @Summary Get user by id
// @Description getHandler returns a single cat fact by id
// @Tags user
// @Produce json
// @Param userID path string true "user id"
// @Router /user/{userID} [get]
// @Success 200 {object} user.User
// @Failure 400 {object} handlerhttp.ErrorResponse
// @Failure 404 {object} handlerhttp.ErrorResponse
func getHandler(
	us *user.Storage,
	gnup handlerhttp.NamedURLParamsGetter,
) handlerhttp.TypedHandler {
	return func(r *http.Request) (*handlerhttp.Response, *handlerhttp.ErrorResponse) {
		userID, errR := gnup(r.Context(), "userID")
		if errR != nil {
			return nil, errR
		}

		id, err := uuid.Parse(userID)
		if err != nil {
			return nil, handlerhttp.ParsingParamError{
				Name:  "userID",
				Value: userID,
			}.ToErrorResponse()
		}

		user, err := us.Get(id)
		if err != nil {
			return nil, handlerhttp.NotFoundError{Designation: "user"}.ToErrorResponse()
		}

		return &handlerhttp.Response{
			Body:           user,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}
