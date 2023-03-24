package middlewares

import "net/http"

func MakeSetRequestPathPrefixMiddleware(prefix string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = prefix + r.URL.Path
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}
