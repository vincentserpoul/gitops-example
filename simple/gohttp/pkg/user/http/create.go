package http

import (
	"gohttp/pkg/user"
	"net/http"

	"github.com/induzo/httpwrapper"
)

// createHandler creates a user
// @Summary creates a user
// @Description createHandler creates a user
// @Tags User
// @Router /user [post]
// @Param user body user.User true "user"
// @Success 201
// @Failure 400 {object} httpwrapper.ErrorResponse
// @Failure 404 {object} httpwrapper.ErrorResponse
func createHandler(us *user.Storage) httpwrapper.TypedHandler {
	return func(r *http.Request) (*httpwrapper.Response, *httpwrapper.ErrorResponse) {
		var user user.User

		if err := httpwrapper.BindBody(r, &user); err != nil {
			return nil, err
		}

		us.Create(&user)

		return &httpwrapper.Response{
			Body:           user,
			HTTPStatusCode: http.StatusCreated,
		}, nil
	}
}
