package happycat

import (
	"archiver/pkg/handler"
	"archiver/pkg/internal/db"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func CreateHandler(
	q db.Querier,
	sugar *zap.SugaredLogger,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params db.SaveHappycatFactParams

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			if err := render.Render(
				w, r,
				handler.ErrRender(
					w, r,
					newErrorResponse(WrongBodyFormatError{Err: err}),
				),
			); err != nil {
				sugar.Errorf("error rendering error: %v", err)
			}

			return
		}

		if err := q.SaveHappycatFact(r.Context(), params); err != nil {
			if err := render.Render(w, r, handler.ErrRender(w, r, newErrorResponse(err))); err != nil {
				sugar.Errorf("error rendering error: %v", err)
			}

			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}
