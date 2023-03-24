package middlewares

import "net/http"

func MakeDedupeResponseHeaders(headers ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, header := range headers {
				w.Header().Set(header, w.Header().Get(header))
			}
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}
