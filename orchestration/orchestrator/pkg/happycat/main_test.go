package happycat

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func testServer(
	t *testing.T,
	serverResponse string,
	timeout time.Duration,
	timeoutFreq int,
) *httptest.Server {
	t.Helper()

	handler := func() func(http.ResponseWriter, *http.Request) {
		count := int64(0)

		return func(w http.ResponseWriter, r *http.Request) {
			localCount := atomic.AddInt64(&count, 1)

			if timeoutFreq != 0 && localCount%int64(timeoutFreq) == 0 {
				time.Sleep(timeout * 2)
			}

			fmt.Fprint(w, serverResponse)
		}
	}

	return httptest.NewServer(
		http.HandlerFunc(handler()))
}
