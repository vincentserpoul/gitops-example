package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestErrorResponse_render(t *testing.T) {
	t.Parallel()

	zl := zerolog.New(io.Discard).With()

	tests := []struct {
		name           string
		accept         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "happy path",
			accept:         "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_code":"test_render","error_msg":"test error user"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", &bytes.Reader{})
			req.Header.Add("Accept", tt.accept)

			nr := httptest.NewRecorder()

			her := &ErrorResponse{
				Error:          fmt.Errorf("test render"),
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "test_render",
				ErrorMsg:       "test error user",
			}

			her.render(zl.Logger(), nr, req)

			resp := nr.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)

				return
			}

			if resp.Header.Get("Content-Type") != tt.accept {
				t.Errorf("expected response header %s, got %s", tt.accept, resp.Header.Get("Content-Type"))

				return
			}

			body, _ := io.ReadAll(resp.Body)
			trimmedBody := strings.TrimSpace(string(body))
			if trimmedBody != tt.expectedBody {
				t.Errorf("expected body\n--%s--\ngot\n--%s--", tt.expectedBody, trimmedBody)

				return
			}
		})
	}
}
