package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/inquizarus/nagg/pkg/httptools"
)

func MakeJWTCopyClaimToHeaderMiddleware(claim, header string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			b64Token, err := httptools.ExtractBearerTokenFromRequest(r)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// We are not "validating" the token and ignore any potential error as other
			// middlewares should have validated it already in terms of authentication
			token, _ := jwt.Parse(b64Token, func(t *jwt.Token) (interface{}, error) {
				return nil, nil
			})

			if claimMap, ok := token.Claims.(jwt.MapClaims); ok {
				var val string
				if val, ok = claimMap[claim].(string); !ok {
					w.WriteHeader(http.StatusBadRequest)
					return

				}
				r.Header.Set(header, val)
			}

			if nil != next {
				next.ServeHTTP(w, r)
			}
		})
	}
}
