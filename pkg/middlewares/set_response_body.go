package middlewares

import (
	"io"
	"net/http"
)

func MakeSetResponseBodyMiddleware(body io.Reader) func(http.Handler) http.Handler {
	data, _ := io.ReadAll(body)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(data)
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}
