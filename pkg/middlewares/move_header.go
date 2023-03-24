package middlewares

import "net/http"

func MakeMoveHeaderMiddleware(source, destination, target string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch target {
			case "request":
				if r.Header.Get(source) != "" {
					r.Header.Set(destination, r.Header.Get(source))
					r.Header.Del(source)
				}
			default:
				if w.Header().Get(source) != "" {
					w.Header().Set(destination, w.Header().Get(source))
					w.Header().Del(source)
				}
			}
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}
