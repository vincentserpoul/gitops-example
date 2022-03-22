package http

import (
	handlerhttp "archiver/pkg/handler/http"
	"archiver/pkg/internal/db"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_getHandler(t *testing.T) {
	t.Parallel()

	gnupGen := func(s string) func(context.Context, string) string {
		return func(_ctx context.Context, key string) string {
			return s
		}
	}

	genericTime, _ := time.Parse("2006-01-02T15:04:05-0700", "2016-07-25T02:22:33+0000")

	tests := []struct {
		name                  string
		existingHappyCatFacts []*db.HappycatFact
		urlParam              string
		httpResponse          *handlerhttp.Response
		httpErrorResponse     *handlerhttp.ErrorResponse
	}{
		{
			name: "happy path",
			existingHappyCatFacts: []*db.HappycatFact{
				{
					ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
					Fact:      "fact",
					CreatedAt: genericTime,
				},
			},
			urlParam: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			httpResponse: &handlerhttp.Response{
				Body: &db.HappycatFact{
					ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
					Fact:      "fact",
					CreatedAt: genericTime,
				},
				HTTPStatusCode: http.StatusOK,
			},
			httpErrorResponse: nil,
		},
		{
			name: "not found with bad querier",
			existingHappyCatFacts: []*db.HappycatFact{
				{
					ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c6"),
					Fact:      "fact",
					CreatedAt: genericTime,
				},
			},
			urlParam:     "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			httpResponse: nil,
			httpErrorResponse: &handlerhttp.ErrorResponse{
				Error:          NoHappyCatFactFoundError{},
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorCode:      "no_cat_found",
				ErrorMsg:       NoHappyCatFactFoundError{}.Error(),
			},
		},
		{
			name:                  "bad url param",
			existingHappyCatFacts: []*db.HappycatFact{},
			urlParam:              "wrong format",
			httpResponse:          nil,
			httpErrorResponse: &handlerhttp.ErrorResponse{
				Error:          WrongIDFormatError{ID: "wrong format", Err: nil},
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "bad_request",
				ErrorMsg:       "wrong format id",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q := newTestQuerier(tt.existingHappyCatFacts)
			req := httptest.NewRequest("GET", "/", nil)

			resp, err := getHandler(q, gnupGen(tt.urlParam))(req)

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

func Benchmark_getHandler(b *testing.B) {
	id := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	genericTime, _ := time.Parse("2006-01-02T15:04:05-0700", "2016-07-25T02:22:33+0000")

	q := newTestQuerier([]*db.HappycatFact{
		{
			ID:        id,
			Fact:      "fact",
			CreatedAt: genericTime,
		},
	})

	req := httptest.NewRequest("GET", "/", nil)

	h := getHandler(
		q,
		func(_ctx context.Context, key string) string { return id.String() },
	)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}
