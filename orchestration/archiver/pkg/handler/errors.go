package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

func ErrRender(_ http.ResponseWriter, r *http.Request, errR *ErrorResponse) *ErrorUser {
	log.Printf("%v\n", errR)

	render.Status(r, errR.HTTPStatusCode)

	return &errR.ErrorUser
}

type ErrorResponse struct {
	Err            error
	HTTPStatusCode int
	ErrorUser
}

type ErrorUser string

func (eu *ErrorUser) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
