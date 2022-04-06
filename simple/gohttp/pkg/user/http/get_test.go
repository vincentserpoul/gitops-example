package http

import (
	"context"
	handlerhttp "gohttp/pkg/handler/http"
	"gohttp/pkg/user"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_getHandler(t *testing.T) {
	t.Parallel()

	gnupGen := func(s string) func(context.Context, string) (string, *handlerhttp.ErrorResponse) {
		return func(_ctx context.Context, key string) (string, *handlerhttp.ErrorResponse) {
			return s, nil
		}
	}

	id := uuid.New()
	createdAt, _ := time.Parse(time.RFC3339, "2022-04-06 03:18:36 +0000 UTC")

	storage := user.NewStorage()
	storage.Create(&user.User{ID: id, CreatedAt: createdAt})

	tests := []struct {
		name              string
		urlParam          string
		httpResponse      *handlerhttp.Response
		httpErrorResponse *handlerhttp.ErrorResponse
	}{
		{
			name:     "happy path",
			urlParam: id.String(),
			httpResponse: &handlerhttp.Response{
				Body: &user.User{
					ID:        id,
					CreatedAt: createdAt,
				},
				HTTPStatusCode: http.StatusOK,
			},
			httpErrorResponse: nil,
		},
		{
			name:         "not found with bad querier",
			urlParam:     "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			httpResponse: nil,
			httpErrorResponse: &handlerhttp.ErrorResponse{
				Error:          user.NoUserFoundError{},
				HTTPStatusCode: http.StatusNotFound,
				ErrorCode:      "no_user_found",
				ErrorMsg:       user.NoUserFoundError{}.Error(),
			},
		},
		{
			name:              "bad url param",
			urlParam:          "wrong format",
			httpResponse:      nil,
			httpErrorResponse: handlerhttp.ParsingParamError{Name: "userID", Value: "wrong format"}.ToErrorResponse(),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", "/", nil)

			resp, err := getHandler(storage, gnupGen(tt.urlParam))(req)

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
	id := uuid.New()
	createdAt, _ := time.Parse(time.RFC3339, "2022-04-06T03:18:36+00:00")

	storage := user.NewStorage()
	storage.Create(&user.User{ID: id, CreatedAt: createdAt})

	req := httptest.NewRequest("GET", "/", nil)

	h := getHandler(
		storage,
		func(
			_ctx context.Context,
			key string,
		) (string, *handlerhttp.ErrorResponse) {
			return id.String(), nil
		},
	)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}
