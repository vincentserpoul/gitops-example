package http

import (
	"bytes"
	"fmt"
	handlerhttp "gohttp/pkg/handler/http"
	"gohttp/pkg/user"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_createHandler(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	createdAt, _ := time.Parse(time.RFC3339, "2022-04-06T03:18:36+00:00")
	storage := user.NewStorage()

	tests := []struct {
		name              string
		payload           string
		httpResponse      *handlerhttp.Response
		httpErrorResponse *handlerhttp.ErrorResponse
	}{
		{
			name:    "happy path",
			payload: fmt.Sprintf(`{"id": "%s","created_at": "%s"}`, id.String(), createdAt.Format(time.RFC3339)),
			httpResponse: &handlerhttp.Response{
				Body:           user.User{ID: id, CreatedAt: createdAt},
				HTTPStatusCode: http.StatusCreated,
			},
			httpErrorResponse: nil,
		},
		{
			name:         "wrong formatted json",
			payload:      "{",
			httpResponse: nil,
			httpErrorResponse: &handlerhttp.ErrorResponse{
				Error:          nil,
				ErrorCode:      "wrong",
				HTTPStatusCode: http.StatusBadRequest,
				ErrorMsg:       "",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "/test/user", bytes.NewReader([]byte(tt.payload)))

			resp, err := createHandler(storage)(req)

			if !reflect.DeepEqual(resp, tt.httpResponse) {
				t.Errorf("expected response %v, got %v (err: %v)", tt.httpResponse, resp, err)

				return
			}

			if err != nil && tt.httpErrorResponse != nil && err.HTTPStatusCode != tt.httpErrorResponse.HTTPStatusCode {
				t.Errorf("expected error %v, got %v", tt.httpErrorResponse, err)

				return
			}
		})
	}
}

func Benchmark_createHandler(b *testing.B) {
	storage := user.NewStorage()
	h := createHandler(storage)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()

		payload := fmt.Sprintf("{\"id\": \"%s\",\"created_at\": \"%s\"}", uuid.New().String(), time.Now().Format(time.RFC3339))

		req := httptest.NewRequest("POST", "/test/user", bytes.NewReader([]byte(payload)))

		b.StartTimer()

		_, _ = h(req)
	}
}
