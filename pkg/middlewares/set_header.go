package middlewares

import "net/http"

func MakeSetHeaderMiddleware(key, value, target string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch target {
			case "request":
				r.Header.Set(key, value)
			default:
				w.Header().Set(key, value)
			}
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}
