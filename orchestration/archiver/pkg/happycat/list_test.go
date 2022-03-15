package happycat

import (
	"archiver/pkg/internal/db"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestListHandler(t *testing.T) {
	t.Parallel()

	sugar := zap.NewNop().Sugar()

	genericTime, _ := time.Parse("2006-01-02T15:04:05-0700", "2016-07-25T02:22:33+0000")

	tests := []struct {
		name                  string
		existingHappyCatFacts []*db.HappycatFact
		wantStatusCode        int
		wantBody              string
	}{
		{
			name: "happy path",
			existingHappyCatFacts: []*db.HappycatFact{
				{
					ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
					Fact:      "fact",
					CreatedAt: genericTime,
				},
				{
					ID:        uuid.MustParse("7ba7b810-9dad-11d1-80b4-00c04fd430c8"),
					Fact:      "fact",
					CreatedAt: genericTime,
				},
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `[{"id":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","fact":"fact","created_at":"2016-07-25T02:22:33Z"},{"id":"7ba7b810-9dad-11d1-80b4-00c04fd430c8","fact":"fact","created_at":"2016-07-25T02:22:33Z"}]`,
		},
		{
			name:                  "not found",
			existingHappyCatFacts: nil,
			wantStatusCode:        http.StatusNotFound,
			wantBody:              `"no happy cat facts have been found"`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q := newTestQuerier(tt.existingHappyCatFacts)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test/happycatfact", nil)

			r := chi.NewRouter()
			AddRoutes("test", r, q, sugar)
			r.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatusCode {
				t.Errorf(
					"handler returned unexpected status: got %d want %d",
					rr.Code, tt.wantStatusCode,
				)

				return
			}

			dataB, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("read error %v", err)

				return
			}

			data := strings.TrimSpace(string(dataB))
			if data != tt.wantBody {
				t.Errorf("expected body \n%s\n, got \n%s\n", tt.wantBody, data)
			}
		})
	}
}

func BenchmarkListHandler(b *testing.B) {
	genericTime, _ := time.Parse("2006-01-02T15:04:05-0700", "2016-07-25T02:22:33+0000")

	sugar := zap.NewNop().Sugar()

	q := newTestQuerier([]*db.HappycatFact{
		{
			ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
			Fact:      "fact",
			CreatedAt: genericTime,
		},
		{
			ID:        uuid.MustParse("7ba7b810-9dad-11d1-80b4-00c04fd430c8"),
			Fact:      "fact",
			CreatedAt: genericTime,
		},
	})

	r := chi.NewRouter()
	AddRoutes("test", r, q, sugar)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()

		req := httptest.NewRequest("GET", "/test/happycatfact", nil)

		b.StartTimer()

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		b.StopTimer()

		res := rr.Result()
		res.Body.Close()

		b.StartTimer()
	}
}
