package middlewares

import "net/http"

func MakeRemoveHeaderMiddleware(source, target string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch target {
			case "request":
				r.Header.Del(source)
			default:
				w.Header().Del(source)
			}
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}
