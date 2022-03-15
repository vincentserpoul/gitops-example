package happycat

import (
	"archiver/pkg/handler"
	"archiver/pkg/internal/db"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Response struct {
	*db.HappycatFact
}

func NewResponse(hcf *db.HappycatFact) *Response {
	return &Response{HappycatFact: hcf}
}

func (hcfr *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetHandler(
	q db.Querier,
	sugar *zap.SugaredLogger,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		happyCatFactID := chi.URLParam(r, "happyCatFactID")

		id, err := uuid.Parse(chi.URLParam(r, "happyCatFactID"))
		if err != nil {
			if err := render.Render(
				w, r,
				handler.ErrRender(
					w, r,
					newErrorResponse(WrongIDFormatError{ID: happyCatFactID, Err: err}),
				),
			); err != nil {
				sugar.Errorf("error rendering error: %v", err)
			}
		}

		hcf, err := q.GetHappycatFact(r.Context(), id)
		if err != nil {
			if err := render.Render(w, r, handler.ErrRender(w, r, newErrorResponse(err))); err != nil {
				sugar.Errorf("error rendering error: %v", err)
			}

			return
		}

		if err := render.Render(w, r, NewResponse(hcf)); err != nil {
			sugar.Errorf("error rendering: %v", err)
		}
	})
}
