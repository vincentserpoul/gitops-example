package http

import (
	handlerhttp "gohttp/pkg/handler/http"
	"gohttp/pkg/user"
	"net/http"
)

// createHandler creates a user
// @Summary creates a user
// @Description createHandler creates a user
// @Tags User
// @Router /user [post]
// @Param user body user.User true "user"
// @Success 201
// @Failure 400 {object} handlerhttp.ErrorResponse
// @Failure 404 {object} handlerhttp.ErrorResponse
func createHandler(us *user.Storage) handlerhttp.TypedHandler {
	return func(r *http.Request) (*handlerhttp.Response, *handlerhttp.ErrorResponse) {
		var user user.User

		if err := handlerhttp.BindBody(r, &user); err != nil {
			return nil, err
		}

		us.Create(&user)

		return &handlerhttp.Response{
			Body:           user,
			HTTPStatusCode: http.StatusCreated,
		}, nil
	}
}
