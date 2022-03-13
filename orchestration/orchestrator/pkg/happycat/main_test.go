package happycat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func testServer(
	t *testing.T,
	serverResponse string,
	timeout time.Duration,
	timeoutFreq int,
) *httptest.Server {
	t.Helper()

	if timeoutFreq == 0 {
		timeoutFreq = 1
	}

	handler := func() func(http.ResponseWriter, *http.Request) {
		count := 0

		return func(w http.ResponseWriter, r *http.Request) {
			count++
			if count%timeoutFreq == 0 && timeout != 0 {
				time.Sleep(timeout * 2)
			}

			fmt.Fprint(w, serverResponse)
		}
	}

	return httptest.NewServer(
		http.HandlerFunc(handler()))
}
