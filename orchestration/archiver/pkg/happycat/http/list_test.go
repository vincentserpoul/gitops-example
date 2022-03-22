package http

import (
	handlerhttp "archiver/pkg/handler/http"
	"archiver/pkg/internal/db"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_listHandler(t *testing.T) {
	t.Parallel()

	genericTime, _ := time.Parse("2006-01-02T15:04:05-0700", "2016-07-25T02:22:33+0000")

	tests := []struct {
		name                  string
		existingHappyCatFacts []*db.HappycatFact
		httpResponse          *handlerhttp.Response
		httpErrorResponse     *handlerhttp.ErrorResponse
	}{
		{
			name: "happy path",
			existingHappyCatFacts: []*db.HappycatFact{
				{
					ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
					Fact:      "fact1",
					CreatedAt: genericTime,
				},
				{
					ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c9"),
					Fact:      "fact2",
					CreatedAt: genericTime,
				},
			},
			httpResponse: &handlerhttp.Response{
				Body: []*db.HappycatFact{
					{
						ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
						Fact:      "fact1",
						CreatedAt: genericTime,
					},
					{
						ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c9"),
						Fact:      "fact2",
						CreatedAt: genericTime,
					},
				},
				HTTPStatusCode: http.StatusOK,
			},
			httpErrorResponse: nil,
		},
		{
			name:                  "not found with bad querier",
			existingHappyCatFacts: []*db.HappycatFact{},
			httpResponse:          nil,
			httpErrorResponse: &handlerhttp.ErrorResponse{
				Error:          NoHappyCatFactFoundError{},
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorCode:      "no_cat_found",
				ErrorMsg:       NoHappyCatFactFoundError{}.Error(),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q := newTestQuerier(tt.existingHappyCatFacts)
			req := httptest.NewRequest("GET", "/", nil)

			resp, err := listHandler(q)(req)

			if tt.httpResponse != nil && !reflect.DeepEqual(resp, tt.httpResponse) {
				t.Errorf("expected response %v, got %v (err: %v)", tt.httpResponse, resp, err)

				return
			}

			if tt.httpErrorResponse != nil {
				if err == nil {
					t.Errorf("expected error but got none")

					return
				}

				if tt.httpErrorResponse != nil && err.HTTPStatusCode != tt.httpErrorResponse.HTTPStatusCode {
					t.Errorf("expected error %d, got %d", tt.httpErrorResponse.HTTPStatusCode, err.HTTPStatusCode)

					return
				}
			}
		})
	}
}

func Benchmark_listHandler(b *testing.B) {
	id := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	genericTime, _ := time.Parse("2006-01-02T15:04:05-0700", "2016-07-25T02:22:33+0000")

	q := newTestQuerier([]*db.HappycatFact{
		{
			ID:        id,
			Fact:      "fact",
			CreatedAt: genericTime,
		},
	})

	h := listHandler(q)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()

		req := httptest.NewRequest("GET", "/", nil)

		b.StartTimer()

		_, _ = h(req)

		b.StopTimer()
	}
}
