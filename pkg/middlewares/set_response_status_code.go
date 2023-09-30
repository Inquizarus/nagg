package middlewares

import (
	"net/http"
)

func MakeSetResponseStatusCodeMiddleware(statusCode int) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(statusCode)
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}
