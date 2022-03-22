package http

import (
	handlerhttp "archiver/pkg/handler/http"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestCreateHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		payload           string
		httpResponse      *handlerhttp.Response
		httpErrorResponse *handlerhttp.ErrorResponse
	}{
		{
			name:    "happy path",
			payload: fmt.Sprintf(`{"id": "%s","fact": "fact"}`, uuid.New().String()),
			httpResponse: &handlerhttp.Response{
				Body:           nil,
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

			q := newTestQuerier(nil)
			req := httptest.NewRequest("POST", "/test/happycatfact", bytes.NewReader([]byte(tt.payload)))

			resp, err := createHandler(q)(req)

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

func BenchmarkCreateHandler(b *testing.B) {
	q := newTestQuerier(nil)
	h := createHandler(q)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()

		payload := fmt.Sprintf("{\"id\": \"%s\",\"fact\": \"fact\"}", uuid.New().String())

		req := httptest.NewRequest("POST", "/test/happycatfact", bytes.NewReader([]byte(payload)))

		b.StartTimer()

		_, _ = h(req)
	}
}
