package middlewares

import "net/http"

func MakeRedirectMiddleware(status int, destination string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			http.Redirect(w, r, destination, status)
		})
	}
}
