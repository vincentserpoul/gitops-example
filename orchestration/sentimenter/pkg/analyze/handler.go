package analyze

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func Handler(analyzer func(string) uint8) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody := struct {
			Text string
		}{}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Printf("error decoding body: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		score := analyzer(reqBody.Text)

		if _, err := w.Write([]byte(strconv.Itoa(int(score)))); err != nil {
			log.Printf("error not written: %v", err)
		}
	})
}
