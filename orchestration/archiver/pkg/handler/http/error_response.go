package http

import (
	"net/http"

	"github.com/rs/zerolog"
)

type ErrorResponse struct {
	Error          error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	ErrorCode      string `json:"error_code"`
	ErrorMsg       string `json:"error_msg"`
}

func (her *ErrorResponse) render(log zerolog.Logger, w http.ResponseWriter, r *http.Request) {
	render(
		log,
		r.Header.Get("Accept"),
		her.HTTPStatusCode,
		her,
		w,
	)
}
