package happycat

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)
	}
}

func testServer(
	t *testing.T,
	serverResponse string,
	timeout time.Duration,
	timeoutFreq int,
) *httptest.Server {
	t.Helper()

	handler := func() func(http.ResponseWriter, *http.Request) {
		count := 0

		return func(w http.ResponseWriter, r *http.Request) {
			count++
			if timeoutFreq != 0 && count%timeoutFreq == 0 {
				time.Sleep(timeout * 2)
			}

			fmt.Fprint(w, serverResponse)
		}
	}

	return httptest.NewServer(
		http.HandlerFunc(handler()))
}
