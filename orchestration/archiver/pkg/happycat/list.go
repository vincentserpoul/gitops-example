package happycat

import (
	"archiver/pkg/handler"
	"archiver/pkg/internal/db"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func NewHappyCatFactListResponse(hcfs []*db.HappycatFact) []render.Renderer {
	list := []render.Renderer{}

	for _, hcf := range hcfs {
		list = append(list, NewResponse(hcf))
	}

	return list
}

func ListHandler(
	q db.Querier,
	sugar *zap.SugaredLogger,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hcfs, err := q.ListHappycatFacts(r.Context())
		if err != nil {
			if err := render.Render(w, r, handler.ErrRender(w, r, newErrorResponse(err))); err != nil {
				sugar.Errorf("error rendering error: %v", err)
			}

			return
		}

		if err := render.RenderList(w, r, NewHappyCatFactListResponse(hcfs)); err != nil {
			sugar.Errorf("error rendering: %v", err)
		}
	})
}
