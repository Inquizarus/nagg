package middlewares

import "net/http"

func MakeRemoveRequestParameterMiddleware(key string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			query.Del(key)
			r.URL.RawQuery = query.Encode()
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}
