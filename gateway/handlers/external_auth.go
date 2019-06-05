package handlers

import (
	"context"
	"net/http"
	"time"
)

// MakeExternalAuthHandler make an authentication proxy handler
func MakeExternalAuthHandler(next http.HandlerFunc, upstreamTimeout time.Duration, upstreamURL string, passBody bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, _ := http.NewRequest(http.MethodGet, upstreamURL, nil)

		deadlineContext, cancel := context.WithTimeout(
			context.Background(),
			upstreamTimeout)

		defer cancel()

		res, err := http.DefaultClient.Do(req.WithContext(deadlineContext))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		if res.StatusCode == http.StatusOK {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(res.StatusCode)
	}
}
