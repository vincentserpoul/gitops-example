package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeBody(r *http.Request, target interface{}) error {
	// nolint: gocritic
	// LATER: add more encodings
	switch r.Header.Get("Content-Type") {
	default:
		if err := json.NewDecoder(r.Body).Decode(target); err != nil {
			return fmt.Errorf("body %s: %w", r.Body, err)
		}
	}

	return nil
}
